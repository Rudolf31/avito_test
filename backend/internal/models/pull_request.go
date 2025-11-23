package models

import "avito_test/internal/models/dto"

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignResponse struct {
	PR           dto.PullRequestDTO `json:"pr"`
	ReplacedByID string             `json:"replaced_by"`
}
