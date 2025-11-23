package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"avito_test/backend/internal/models"
	"avito_test/backend/internal/models/dto"
	"avito_test/backend/internal/services"
)

// RegisterRoutes registers all HTTP routes and handlers using external service interfaces.
// Handlers now use DTOs/models for requests and responses to keep gin isolated from internals.
func RegisterRoutes(router *gin.Engine, teamService services.TeamService, userService services.UserService, prService services.PullRequestService) {

	var teamSvc services.ExternalTeamService = teamService
	var userSvc services.ExternalUserService = userService
	var prSvc services.ExternalPullRequestService = prService

	// Teams
	router.POST("/team/add", func(c *gin.Context) {
		var req dto.TeamDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": err.Error()}})
			return
		}

		resp, err := teamSvc.APIAddTeam(req)
		if err != nil {
			if strings.Contains(err.Error(), "team already exists") {
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "TEAM_EXISTS", "message": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
			return
		}

		// models.TeamAddResponse already has json tags ("team"), so we can return it directly.
		c.JSON(http.StatusCreated, resp)
	})

	router.GET("/team/get", func(c *gin.Context) {
		teamName := c.Query("team_name")
		if teamName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": "team_name required"}})
			return
		}

		resp, err := teamSvc.APIGetTeam(teamName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
			return
		}

		// dto.TeamDTO has json tags team_name/members
		c.JSON(http.StatusOK, resp)
	})

	// Users
	router.POST("/users/setIsActive", func(c *gin.Context) {
		var req models.SetIsActiveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": err.Error()}})
			return
		}

		resp, err := userSvc.APISetIsActive(req)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
			return
		}

		// models.SetIsActiveResponse has json tag user
		c.JSON(http.StatusOK, resp)
	})

	// PullRequests
	router.POST("/pullRequest/create", func(c *gin.Context) {
		var req dto.PullRequestShortDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": err.Error()}})
			return
		}

		pr, err := prSvc.APICreatePullRequest(req)
		if err != nil {
			if strings.Contains(err.Error(), "PR id already exists") {
				c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "PR_EXISTS", "message": err.Error()}})
				return
			}
			if strings.Contains(err.Error(), "author not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
			return
		}

		// Service returns dto.PullRequestDTO; wrap it under "pr" to match OpenAPI.
		c.JSON(http.StatusCreated, gin.H{"pr": pr})
	})

	router.POST("/pullRequest/merge", func(c *gin.Context) {
		var req struct {
			PullRequestID string `json:"pull_request_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": err.Error()}})
			return
		}

		if err := prSvc.APIMergePullRequest(req.PullRequestID); err != nil {
			if strings.Contains(err.Error(), "pr not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
			return
		}

		// No full PR DTO is returned by the external merge API; return minimal merged info matching DTO fields.
		c.JSON(http.StatusOK, gin.H{"pr": gin.H{"pull_request_id": req.PullRequestID, "status": "MERGED", "merged_at": time.Now().UTC().Format(time.RFC3339)}})
	})

	router.POST("/pullRequest/reassign", func(c *gin.Context) {
		var req models.ReassignRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": err.Error()}})
			return
		}

		resp, err := prSvc.APIReassignPullRequest(req)
		if err != nil {
			if strings.Contains(err.Error(), "pr not found") || strings.Contains(err.Error(), "user not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
				return
			}
			if strings.Contains(err.Error(), "cannot reassign on merged PR") {
				c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "PR_MERGED", "message": err.Error()}})
				return
			}
			if strings.Contains(err.Error(), "reviewer is not assigned to this PR") {
				c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "NOT_ASSIGNED", "message": err.Error()}})
				return
			}
			if strings.Contains(err.Error(), "no active replacement candidate in team") {
				c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "NO_CANDIDATE", "message": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
			return
		}

		// models.ReassignResponse already has json tags (pr, replaced_by)
		c.JSON(http.StatusOK, resp)
	})

	router.GET("/users/getReview", func(c *gin.Context) {
		userID := c.Query("user_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "INVALID_PAYLOAD", "message": "user_id required"}})
			return
		}

		resp, err := userSvc.APIGetReview(userID)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": err.Error()}})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
			return
		}

		// models.GetReviewResponse uses dto.PullRequestDTO for items; return directly.
		c.JSON(http.StatusOK, resp)
	})
}
