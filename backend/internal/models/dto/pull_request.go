package dto

import "avito_test/internal/models/enum"

type PullRequestDTO struct {
	ID                string        `json:"pull_request_id"`
	PullRequestName   string        `json:"pull_request_name"`
	AuthorID          string        `json:"author_id"`
	Status            enum.PrStatus `json:"status"`
	AssignedReviewers []string      `json:"assigned_reviewers,omitempty"`
	CreatedAt         string        `json:"created_at"`
	MergedAt          string        `json:"merged_at,omitempty"`
}

type PullRequestShortDTO struct {
	ID              string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}
