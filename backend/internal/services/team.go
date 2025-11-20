package services

import (
	"avito_test/backend/internal/db"
)

type InternalTeamService interface {
	GetTeamByID(teamID string) (*db.Team, error)
	CreateTeam(name string) (*db.Team, error)
	AddUserToTeamByID(teamID string, userID string) (*db.Team, error)
	UpdateTeam(teamID string, name string) (*db.Team, error)
}

type ExternalTeamService interface {
	// TODO: Create DTO
	APIAddTeam(name string) (*db.Team, error)
	APIGetTeam(teamID string) (*db.Team, error)
}

type TeamService interface {
	InternalTeamService
	ExternalTeamService
}
