package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/config"
	"github.com/seta-namnv-6798/go-apis/models"
)

// CreateFolderRequest represents the request structure for creating a folder
type CreateFolderRequest struct {
	Name    string `json:"name" binding:"required"`
	OwnerID uint   `json:"ownerId" binding:"required"`
}

// UpdateFolderRequest represents the request structure for updating a folder
type UpdateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}

// ShareFolderRequest represents the request for sharing a folder
type ShareFolderRequest struct {
	UserID uint   `json:"userId" binding:"required"`
	Access string `json:"access" binding:"required,oneof=read write"`
}

// CreateFolder creates a new folder
func CreateFolder(c *gin.Context) {
	var req CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if owner exists
	var owner models.User
	if err := config.DB.First(&owner, req.OwnerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Owner not found"})
		return
	}

	// Create the folder
	folder := models.Folder{
		Name:    req.Name,
		OwnerID: req.OwnerID,
	}

	if err := config.DB.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}

	// Load the folder with owner information
	config.DB.Preload("Owner").First(&folder, folder.FolderID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Folder created successfully",
		"folder":  folder,
	})
}

// GetFolder retrieves folder details
func GetFolder(c *gin.Context) {
	folderIDStr := c.Param("folderId")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var folder models.Folder
	if err := config.DB.Preload("Owner").First(&folder, folderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"folder": folder,
	})
}

// UpdateFolder updates folder information
func UpdateFolder(c *gin.Context) {
	folderIDStr := c.Param("folderId")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var folder models.Folder
	if err := config.DB.First(&folder, folderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Update folder
	folder.Name = req.Name
	if err := config.DB.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}

	// Load updated folder with owner
	config.DB.Preload("Owner").First(&folder, folder.FolderID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Folder updated successfully",
		"folder":  folder,
	})
}

// DeleteFolder deletes a folder and all its notes
func DeleteFolder(c *gin.Context) {
	folderIDStr := c.Param("folderId")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var folder models.Folder
	if err := tx.First(&folder, folderID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Delete all note shares for notes in this folder
	if err := tx.Where("note_id IN (SELECT note_id FROM notes WHERE folder_id = ?)", folderID).Delete(&models.NoteShare{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note shares"})
		return
	}

	// Delete all notes in the folder
	if err := tx.Where("folder_id = ?", folderID).Delete(&models.Note{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notes"})
		return
	}

	// Delete folder shares
	if err := tx.Where("folder_id = ?", folderID).Delete(&models.FolderShare{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder shares"})
		return
	}

	// Delete the folder
	if err := tx.Delete(&folder).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder and its contents deleted successfully"})
}

// ShareFolder shares a folder with a user
func ShareFolder(c *gin.Context) {
	folderIDStr := c.Param("folderId")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req ShareFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if folder exists
	var folder models.Folder
	if err := config.DB.First(&folder, folderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if folder is already shared with this user
	var existingShare models.FolderShare
	if err := config.DB.Where("folder_id = ? AND user_id = ?", folderID, req.UserID).First(&existingShare).Error; err == nil {
		// Update existing share
		existingShare.Access = req.Access
		if err := config.DB.Save(&existingShare).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder share"})
			return
		}
		config.DB.Preload("User").Preload("Folder").First(&existingShare, existingShare.ID)
		c.JSON(http.StatusOK, gin.H{
			"message": "Folder share updated successfully",
			"share":   existingShare,
		})
		return
	}

	// Create new share
	folderShare := models.FolderShare{
		FolderID: uint(folderID),
		UserID:   req.UserID,
		Access:   req.Access,
	}

	if err := config.DB.Create(&folderShare).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share folder"})
		return
	}

	// Load share with relationships
	config.DB.Preload("User").Preload("Folder").First(&folderShare, folderShare.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Folder shared successfully",
		"share":   folderShare,
	})
}

// RevokeFolderShare revokes folder sharing for a user
func RevokeFolderShare(c *gin.Context) {
	folderIDStr := c.Param("folderId")
	userIDStr := c.Param("userId")

	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find and delete the folder share
	var folderShare models.FolderShare
	if err := config.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).First(&folderShare).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder share not found"})
		return
	}

	if err := config.DB.Delete(&folderShare).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke folder share"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder share revoked successfully"})
}
