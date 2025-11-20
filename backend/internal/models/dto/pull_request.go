package dto

type PullRequestDTO struct {
	ID              string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	// TODO: Change to enum
	Status             string `json:"status"`
	AssignedReviewerID string `json:"assigned_reviewer_id,omitempty"`
	CreatedAt          string `json:"created_at"`
	MergedAt           string `json:"merged_at,omitempty"`
}

type PullRequestShortDTO struct {
	ID              string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	// TODO: Change to enum
	Status string `json:"status"`
}
