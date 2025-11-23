package dto

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamMemberDTO struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
