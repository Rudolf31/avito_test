package services

import (
	"avito_test/internal/db"
	"avito_test/internal/models"
	"avito_test/internal/models/dto"
	"avito_test/internal/models/enum"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InternalUserService interface {
	GetUserByID(userID string) (*db.User, error)
	CreateUser(username, id string) (*db.User, error)
	SetIsActive(userID string, isActive bool) error
}

type ExternalUserService interface {
	APISetIsActive(request models.SetIsActiveRequest) (*models.SetIsActiveResponse, error)
	APIGetReview(userID string) (*models.GetReviewResponse, error)
}

type UserService interface {
	InternalUserService
	ExternalUserService
}

type userServiceImpl struct {
	pool *pgxpool.Pool
}

// APIGetReview implements UserService.
func (u *userServiceImpl) APIGetReview(userID string) (*models.GetReviewResponse, error) {
	q := db.New(u.pool)
	ctx := context.TODO()

	// 1. Проверяем, что пользователь существует
	_, err := q.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 2. Берём все записи review где user_id = reviewer
	reviews, err := q.GetReviewsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	pullRequests := make([]dto.PullRequestDTO, 0, len(reviews))

	// 3. Для каждого review отдаём информацию о PR
	for _, r := range reviews {
		pr, err := q.GetPullRequestByID(ctx, r.PullRequestID)
		if err != nil {
			return nil, err
		}

		pullRequests = append(pullRequests, dto.PullRequestDTO{
			ID:              pr.ID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          enum.PrStatus(pr.Status),
		})
	}

	// 4. Собираем финальный ответ
	return &models.GetReviewResponse{
		UserID:       userID,
		PullRequests: pullRequests,
	}, nil
}

// APISetIsActive implements UserService.
func (u *userServiceImpl) APISetIsActive(request models.SetIsActiveRequest) (*models.SetIsActiveResponse, error) {
	q := db.New(u.pool)

	ctx := context.TODO()

	err := u.SetIsActive(request.UserID, request.IsActive)
	if err != nil {
		return nil, err
	}

	result, err := q.GetUserWithTeamNameByID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	response := &models.SetIsActiveResponse{
		User: dto.UserDTO{
			ID:       result.ID,
			Username: result.Username,
			IsActive: result.IsActive,
			TeamName: result.TeamName,
		},
	}
	return response, nil
}

// CreateUser implements UserService.
func (u *userServiceImpl) CreateUser(username, id string) (*db.User, error) {

	q := db.New(u.pool)

	ctx := context.TODO()

	user := db.CreateUserParams{
		ID:       id,
		Username: username,
		IsActive: true,
	}

	newUser, err := q.CreateUser(ctx, user)

	if err != nil {
		//TODO: write specific error handling
		return nil, err
	}

	return &newUser, nil
}

// GetUserByID implements UserService.
func (u *userServiceImpl) GetUserByID(userID string) (*db.User, error) {
	q := db.New(u.pool)

	ctx := context.TODO()

	user, err := q.GetUserByID(ctx, userID)

	if err != nil {
		//TODO: write specific error handling
		return nil, err
	}

	return &user, nil
}

// SetIsActive implements UserService.
func (u *userServiceImpl) SetIsActive(userID string, isActive bool) error {
	q := db.New(u.pool)

	ctx := context.TODO()

	err := q.SetIsActive(ctx, db.SetIsActiveParams{
		ID:       userID,
		IsActive: isActive,
	})

	if err != nil {
		//TODO: write specific error handling
		return err
	}

	return nil
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return &userServiceImpl{
		pool: pool,
	}
}
