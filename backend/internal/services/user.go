package services

import "avito_test/backend/internal/db"

type InternalUserService interface {
	GetUserByID(userID string) (*db.User, error)
	CreateUser(username string) (*db.User, error)
	SetIsActive(userID string, isActive bool) error
}

type ExternalUserService interface {
	// TODO: Create DTO
	APISetIsActive
	APIGetReview
}

type UserService interface {
	InternalUserService
	ExternalUserService
}
