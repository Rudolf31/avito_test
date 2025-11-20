package dto

type UserDTO struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	TeamID   string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
