package services

import

type InternalPullRequestService interface {
	CreatePullRequest(PullRequestID string, PullRequestName string, authorID string) (*db.PullRequest, error)
	MergePullRequest(PullRequestID string) error
	ReassignPullRequest(PullRequestID string, oldReviewerID string) error
}

type ExternalPullRequestService interface {
	// TODO: Create DTO
	APICreatePullRequest
	APIMergePullRequest(PullRequestID string) error
	APIReassignPullRequest
}

type PullRequestService interface {
	InternalPullRequestService
	ExternalPullRequestService
}