package services

import (
	"avito_test/backend/internal/db"
	"avito_test/backend/internal/models"
	"avito_test/backend/internal/models/dto"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InternalTeamService interface {
	GetTeamByID(teamID string) (*db.Team, error)
	CreateTeam(name string) (*db.Team, error)
}

type ExternalTeamService interface {
	APIAddTeam(request dto.TeamDTO) (*models.TeamAddResponse, error)
	APIGetTeam(teamName string) (*dto.TeamDTO, error)
}

type TeamService interface {
	InternalTeamService
	ExternalTeamService
}

type teamServiceImpl struct {
	pool *pgxpool.Pool
}

// APIAddTeam implements TeamService.
func (t *teamServiceImpl) APIAddTeam(request dto.TeamDTO) (*models.TeamAddResponse, error) {

	c := context.TODO()

	tx, err := t.pool.Begin(c)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)
	q := db.New(tx)
	qtx := q.WithTx(tx)

	team, err := qtx.GetTeamByName(c, request.TeamName)

	if err == nil {
		return nil, errors.New("team already exists")
	}

	if errors.Is(err, pgx.ErrNoRows) {
		team, err = qtx.CreateTeam(c, request.TeamName)
		if err != nil {
			return nil, err
		}
	}

	for _, member := range request.Members {
		_, err := qtx.UpsertUser(c, db.UpsertUserParams{
			ID:       member.ID,
			Username: member.Username,
			TeamID:   team.ID,
			IsActive: member.IsActive,
		})
		if err != nil {
			return nil, err
		}
	}

	return &models.TeamAddResponse{
		Team: dto.TeamDTO{
			TeamName: team.TeamName,
			Members:  request.Members,
		},
	}, tx.Commit(c)
}

// APIGetTeam implements TeamService.
func (t *teamServiceImpl) APIGetTeam(teamName string) (*dto.TeamDTO, error) {
	c := context.TODO()

	q := db.New(t.pool)

	team, err := q.GetTeamByName(c, teamName)
	if err != nil {
		return nil, err
	}

	users, err := q.GetUsersByTeamID(c, team.ID)
	if err != nil {
		return nil, err
	}

	var members []dto.TeamMemberDTO
	for _, user := range users {
		members = append(members, dto.TeamMemberDTO{
			ID:       user.ID,
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	return &dto.TeamDTO{
		TeamName: team.TeamName,
		Members:  members,
	}, nil
}

// CreateTeam implements TeamService.
func (t *teamServiceImpl) CreateTeam(name string) (*db.Team, error) {
	q := db.New(t.pool)

	ctx := context.TODO()

	team, err := q.CreateTeam(ctx, name)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// GetTeamByID implements TeamService.
func (t *teamServiceImpl) GetTeamByID(teamID string) (*db.Team, error) {

	q := db.New(t.pool)

	ctx := context.TODO()

	team, err := q.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func NewTeamService(pool *pgxpool.Pool) TeamService {
	return &teamServiceImpl{
		pool: pool,
	}
}
