package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seta-namnv-6798/go-apis/config"
	"github.com/seta-namnv-6798/go-apis/models"
)

// CreateTeamRequest represents the request structure for creating a team
type CreateTeamRequest struct {
	TeamName string               `json:"teamName" binding:"required"`
	Managers []TeamManagerRequest `json:"managers"`
	Members  []TeamMemberRequest  `json:"members"`
}

type TeamManagerRequest struct {
	ManagerID   string `json:"managerId" binding:"required"`
	ManagerName string `json:"managerName" binding:"required"`
}

type TeamMemberRequest struct {
	MemberID   string `json:"memberId" binding:"required"`
	MemberName string `json:"memberName" binding:"required"`
}

// AddMemberRequest represents the request for adding a member to a team
type AddMemberRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// AddManagerRequest represents the request for adding a manager to a team
type AddManagerRequest struct {
	UserID uint `json:"userId" binding:"required"`
}

// CreateTeam creates a new team with managers and members
func CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the team
	team := models.Team{
		TeamName: req.TeamName,
	}

	if err := tx.Create(&team).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	// Add managers
	for _, managerReq := range req.Managers {
		userID, err := strconv.ParseUint(managerReq.ManagerID, 10, 32)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
			return
		}

		// Check if user exists
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Manager not found"})
			return
		}

		teamManager := models.TeamManager{
			UserID: uint(userID),
			TeamID: team.TeamID,
		}

		if err := tx.Create(&teamManager).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add manager"})
			return
		}
	}

	// Add members
	for _, memberReq := range req.Members {
		userID, err := strconv.ParseUint(memberReq.MemberID, 10, 32)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
			return
		}

		// Check if user exists
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
			return
		}

		teamMember := models.TeamMember{
			UserID: uint(userID),
			TeamID: team.TeamID,
		}

		if err := tx.Create(&teamMember).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Return the created team with its relationships
	var createdTeam models.Team
	config.DB.First(&createdTeam, team.TeamID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Team created successfully",
		"team":    createdTeam,
	})
}

// AddMemberToTeam adds a member to an existing team
func AddMemberToTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if team exists
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already a member
	var existingMember models.TeamMember
	if err := config.DB.Where("user_id = ? AND team_id = ?", req.UserID, teamID).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a member of this team"})
		return
	}

	// Add the member
	teamMember := models.TeamMember{
		UserID: req.UserID,
		TeamID: uint(teamID),
	}

	if err := config.DB.Create(&teamMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Member added successfully",
		"member":  teamMember,
	})
}

// RemoveMemberFromTeam removes a member from a team
func RemoveMemberFromTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	memberIDStr := c.Param("memberId")

	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid member ID"})
		return
	}

	// Check if team exists
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Find and delete the team member relationship
	var teamMember models.TeamMember
	if err := config.DB.Where("user_id = ? AND team_id = ?", memberID, teamID).First(&teamMember).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found in team"})
		return
	}

	if err := config.DB.Delete(&teamMember).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// AddManagerToTeam adds a manager to an existing team
func AddManagerToTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req AddManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if team exists
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Check if user exists
	var user models.User
	if err := config.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if user is already a manager
	var existingManager models.TeamManager
	if err := config.DB.Where("user_id = ? AND team_id = ?", req.UserID, teamID).First(&existingManager).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User is already a manager of this team"})
		return
	}

	// Add the manager
	teamManager := models.TeamManager{
		UserID: req.UserID,
		TeamID: uint(teamID),
	}

	if err := config.DB.Create(&teamManager).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add manager"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Manager added successfully",
		"manager": teamManager,
	})
}

// RemoveManagerFromTeam removes a manager from a team
func RemoveManagerFromTeam(c *gin.Context) {
	teamIDStr := c.Param("teamId")
	managerIDStr := c.Param("managerId")

	teamID, err := strconv.ParseUint(teamIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	managerID, err := strconv.ParseUint(managerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid manager ID"})
		return
	}

	// Check if team exists
	var team models.Team
	if err := config.DB.First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	// Find and delete the team manager relationship
	var teamManager models.TeamManager
	if err := config.DB.Where("user_id = ? AND team_id = ?", managerID, teamID).First(&teamManager).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manager not found in team"})
		return
	}

	if err := config.DB.Delete(&teamManager).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove manager"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Manager removed successfully"})
}
