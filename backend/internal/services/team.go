package services

import (
	"avito_test/backend/internal/db"
	"avito_test/backend/internal/models"
	"avito_test/backend/internal/models/dto"
)

type InternalTeamService interface {
	GetTeamByID(teamID string) (*db.Team, error)
	CreateTeam(name string) (*db.Team, error)
	AddUserToTeamByID(teamID string, userID string) (*db.Team, error)
	UpdateTeam(teamID string, name string) (*db.Team, error)
}

type ExternalTeamService interface {
	APIAddTeam(request dto.TeamDTO) (*models.TeamAddResponse, error)
	APIGetTeam(teamName string) (*dto.TeamDTO, error)
}

type TeamService interface {
	InternalTeamService
	ExternalTeamService
}
