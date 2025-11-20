package dto

type TeamDTO struct {
	ID       string `json:"id"`
	TeamName string `json:"team_name"`
	TeamMemberDTO
}

type TeamMemberDTO struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
