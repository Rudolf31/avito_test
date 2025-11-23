package services

import (
	"avito_test/backend/internal/db"
	"avito_test/backend/internal/models"
)

type InternalUserService interface {
	GetUserByID(userID string) (*db.User, error)
	CreateUser(username string) (*db.User, error)
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
