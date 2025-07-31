package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/config"
	"github.com/seta-namnv-6798/go-apis/models"
)

// AssetResponse represents the structure for asset responses
type AssetResponse struct {
	Folders []FolderWithAccess `json:"folders"`
	Notes   []NoteWithAccess   `json:"notes"`
}

type FolderWithAccess struct {
	models.Folder
	AccessType string `json:"accessType"` // "owner", "read", "write"
}

type NoteWithAccess struct {
	models.Note
	AccessType string `json:"accessType"` // "owner", "read", "write"
}

// GetTeamAssets retrieves all assets that team members own or can access (Manager-only)
func GetTeamAssets(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	// Check if team exists
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Get all team members
	var teamMembers []models.TeamMember
	if err := config.DB.Where("team_id = ?", teamID).Find(&teamMembers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get team members"})
		return
	}

	// Extract user IDs
	var userIDs []uint
	for _, member := range teamMembers {
		userIDs = append(userIDs, member.UserID)
	}

	if len(userIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"teamId": teamID,
			"assets": AssetResponse{
				Folders: []FolderWithAccess{},
				Notes:   []NoteWithAccess{},
			},
		})
		return
	}

	var foldersWithAccess []FolderWithAccess
	var notesWithAccess []NoteWithAccess

	// Get folders owned by team members
	var ownedFolders []models.Folder
	config.DB.Preload("Owner").Where("owner_id IN ?", userIDs).Find(&ownedFolders)
	for _, folder := range ownedFolders {
		foldersWithAccess = append(foldersWithAccess, FolderWithAccess{
			Folder:     folder,
			AccessType: "owner",
		})
	}

	// Get folders shared with team members
	var folderShares []models.FolderShare
	config.DB.Preload("Folder").Preload("Folder.Owner").Where("user_id IN ?", userIDs).Find(&folderShares)
	for _, share := range folderShares {
		foldersWithAccess = append(foldersWithAccess, FolderWithAccess{
			Folder:     share.Folder,
			AccessType: share.Access,
		})
	}

	// Get notes owned by team members
	var ownedNotes []models.Note
	config.DB.Preload("Owner").Preload("Folder").Where("owner_id IN ?", userIDs).Find(&ownedNotes)
	for _, note := range ownedNotes {
		notesWithAccess = append(notesWithAccess, NoteWithAccess{
			Note:       note,
			AccessType: "owner",
		})
	}

	// Get notes shared with team members
	var noteShares []models.NoteShare
	config.DB.Preload("Note").Preload("Note.Owner").Preload("Note.Folder").Where("user_id IN ?", userIDs).Find(&noteShares)
	for _, share := range noteShares {
		notesWithAccess = append(notesWithAccess, NoteWithAccess{
			Note:       share.Note,
			AccessType: share.Access,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"teamId": teamID,
		"assets": AssetResponse{
			Folders: foldersWithAccess,
			Notes:   notesWithAccess,
		},
	})
}

// GetUserAssets retrieves all assets owned by or shared with a user
func GetUserAssets(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var foldersWithAccess []FolderWithAccess
	var notesWithAccess []NoteWithAccess

	// Get folders owned by the user
	var ownedFolders []models.Folder
	config.DB.Preload("Owner").Where("owner_id = ?", userID).Find(&ownedFolders)
	for _, folder := range ownedFolders {
		foldersWithAccess = append(foldersWithAccess, FolderWithAccess{
			Folder:     folder,
			AccessType: "owner",
		})
	}

	// Get folders shared with the user
	var folderShares []models.FolderShare
	config.DB.Preload("Folder").Preload("Folder.Owner").Where("user_id = ?", userID).Find(&folderShares)
	for _, share := range folderShares {
		foldersWithAccess = append(foldersWithAccess, FolderWithAccess{
			Folder:     share.Folder,
			AccessType: share.Access,
		})
	}

	// Get notes owned by the user
	var ownedNotes []models.Note
	config.DB.Preload("Owner").Preload("Folder").Where("owner_id = ?", userID).Find(&ownedNotes)
	for _, note := range ownedNotes {
		notesWithAccess = append(notesWithAccess, NoteWithAccess{
			Note:       note,
			AccessType: "owner",
		})
	}

	// Get notes shared with the user
	var noteShares []models.NoteShare
	config.DB.Preload("Note").Preload("Note.Owner").Preload("Note.Folder").Where("user_id = ?", userID).Find(&noteShares)
	for _, share := range noteShares {
		notesWithAccess = append(notesWithAccess, NoteWithAccess{
			Note:       share.Note,
			AccessType: share.Access,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"user":   user,
		"assets": AssetResponse{
			Folders: foldersWithAccess,
			Notes:   notesWithAccess,
		},
	})
}
