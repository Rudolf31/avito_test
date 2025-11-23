package models

import "avito_test/backend/internal/models/dto"

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	User dto.UserDTO `json:"user"`
}

type GetReviewResponse struct {
	UserID       string               `json:"user_id"`
	PullRequests []dto.PullRequestDTO `json:"pull_requests"`
}
