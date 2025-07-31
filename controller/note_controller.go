package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/config"
	"github.com/seta-namnv-6798/go-apis/models"
)

// CreateNoteRequest represents the request structure for creating a note
type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Body    string `json:"body"`
	OwnerID uint   `json:"ownerId" binding:"required"`
}

// UpdateNoteRequest represents the request structure for updating a note
type UpdateNoteRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body"`
}

// ShareNoteRequest represents the request for sharing a note
type ShareNoteRequest struct {
	UserID uint   `json:"userId" binding:"required"`
	Access string `json:"access" binding:"required,oneof=read write"`
}

// CreateNote creates a new note inside a folder
func CreateNote(c *gin.Context) {
	folderIDStr := c.Param("folderId")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	var req CreateNoteRequest
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

	// Check if owner exists
	var owner models.User
	if err := config.DB.First(&owner, req.OwnerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Owner not found"})
		return
	}

	// Create the note
	note := models.Note{
		Title:    req.Title,
		Body:     req.Body,
		FolderID: uint(folderID),
		OwnerID:  req.OwnerID,
	}

	if err := config.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	// Load the note with relationships
	config.DB.Preload("Owner").Preload("Folder").First(&note, note.NoteID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Note created successfully",
		"note":    note,
	})
}

// GetNote retrieves note details
func GetNote(c *gin.Context) {
	noteIDStr := c.Param("noteId")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var note models.Note
	if err := config.DB.Preload("Owner").Preload("Folder").First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

// UpdateNote updates note information
func UpdateNote(c *gin.Context) {
	noteIDStr := c.Param("noteId")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Update note
	note.Title = req.Title
	note.Body = req.Body
	if err := config.DB.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	// Load updated note with relationships
	config.DB.Preload("Owner").Preload("Folder").First(&note, note.NoteID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Note updated successfully",
		"note":    note,
	})
}

// DeleteNote deletes a note
func DeleteNote(c *gin.Context) {
	noteIDStr := c.Param("noteId")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	// Start transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var note models.Note
	if err := tx.First(&note, noteID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Delete note shares
	if err := tx.Where("note_id = ?", noteID).Delete(&models.NoteShare{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note shares"})
		return
	}

	// Delete the note
	if err := tx.Delete(&note).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// ShareNote shares a note with a user
func ShareNote(c *gin.Context) {
	noteIDStr := c.Param("noteId")
	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req ShareNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if note exists
	var note models.Note
	if err := config.DB.First(&note, noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if note is already shared with this user
	var existingShare models.NoteShare
	if err := config.DB.Where("note_id = ? AND user_id = ?", noteID, req.UserID).First(&existingShare).Error; err == nil {
		// Update existing share
		existingShare.Access = req.Access
		if err := config.DB.Save(&existingShare).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note share"})
			return
		}
		config.DB.Preload("User").Preload("Note").First(&existingShare, existingShare.ID)
		c.JSON(http.StatusOK, gin.H{
			"message": "Note share updated successfully",
			"share":   existingShare,
		})
		return
	}

	// Create new share
	noteShare := models.NoteShare{
		NoteID: uint(noteID),
		UserID: req.UserID,
		Access: req.Access,
	}

	if err := config.DB.Create(&noteShare).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share note"})
		return
	}

	// Load share with relationships
	config.DB.Preload("User").Preload("Note").First(&noteShare, noteShare.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Note shared successfully",
		"share":   noteShare,
	})
}

// RevokeNoteShare revokes note sharing for a user
func RevokeNoteShare(c *gin.Context) {
	noteIDStr := c.Param("noteId")
	userIDStr := c.Param("userId")

	noteID, err := strconv.ParseUint(noteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find and delete the note share
	var noteShare models.NoteShare
	if err := config.DB.Where("note_id = ? AND user_id = ?", noteID, userID).First(&noteShare).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note share not found"})
		return
	}

	if err := config.DB.Delete(&noteShare).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke note share"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note share revoked successfully"})
}
