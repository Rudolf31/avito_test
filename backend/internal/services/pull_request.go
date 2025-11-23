package services

import (
	"avito_test/backend/internal/db"
	"avito_test/backend/internal/models"
	"avito_test/backend/internal/models/dto"
)

type InternalPullRequestService interface {
	CreatePullRequest(pullRequestID string, pullRequestName string, authorID string) (*db.PullRequest, error)
	MergePullRequest(pullRequestID string) error
	ReassignPullRequest(pullRequestID string, oldReviewerID string) error
}

type ExternalPullRequestService interface {
	APICreatePullRequest(request dto.PullRequestShortDTO) (*dto.PullRequestDTO, error)
	APIMergePullRequest(pullRequestID string) error
	APIReassignPullRequest(request models.ReassignRequest) (*models.ReassignResponse, error)
}

type PullRequestService interface {
	InternalPullRequestService
	ExternalPullRequestService
}
