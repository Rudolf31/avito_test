package services

import (
	"avito_test/internal/db"
	"avito_test/internal/models"
	"avito_test/internal/models/dto"
	"avito_test/internal/models/enum"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InternalPullRequestService interface {
	CreatePullRequest(pullRequestID string, pullRequestName string, authorID string) (*db.PullRequest, error)
	MergePullRequest(pullRequestID string) error
	ReassignPullRequest(pullRequestID string, oldReviewerID string) error
}

type ExternalPullRequestService interface {
	APICreatePullRequest(request dto.PullRequestShortDTO) (*dto.PullRequestDTO, error)
	APIMergePullRequest(pullRequestID string) error
	APIReassignPullRequest(request models.ReassignRequest) (*models.ReassignResponse, error)
}

type PullRequestService interface {
	InternalPullRequestService
	ExternalPullRequestService
}

type pullRequestServiceImpl struct {
	pool *pgxpool.Pool
}

// APICreatePullRequest implements PullRequestService.
func (p *pullRequestServiceImpl) APICreatePullRequest(request dto.PullRequestShortDTO) (*dto.PullRequestDTO, error) {
	c := context.TODO()

	tx, err := p.pool.Begin(c)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	q := db.New(tx)
	qtx := q.WithTx(tx)

	if _, err := qtx.GetPullRequestByID(c, request.ID); err == nil {
		return nil, errors.New("PR id already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	author, err := qtx.GetUserByID(c, request.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	users, err := qtx.GetUsersByTeamID(c, author.TeamID)
	if err != nil {
		return nil, err
	}

	var assigned []db.User
	for _, u := range users {
		if u.ID == author.ID {
			continue
		}
		if !u.IsActive {
			continue
		}
		assigned = append(assigned, u)
		if len(assigned) == 2 {
			break
		}
	}

	pr, err := qtx.CreatePullRequest(c, db.CreatePullRequestParams{
		ID:              request.ID,
		PullRequestName: request.PullRequestName,
		AuthorID:        request.AuthorID,
		Status:          0,
	})
	if err != nil {
		return nil, err
	}

	var assignedIDs []string
	for _, a := range assigned {
		reviewID := uuid.New().String()
		_, err := qtx.CreateReview(c, db.CreateReviewParams{
			ID:            reviewID,
			UserID:        a.ID,
			PullRequestID: pr.ID,
			Reviewed:      false,
		})
		if err != nil {
			return nil, err
		}
		assignedIDs = append(assignedIDs, a.ID)
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	resp := &dto.PullRequestDTO{
		ID:                pr.ID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            enum.PrStatus(pr.Status),
		AssignedReviewers: assignedIDs,
	}
	return resp, nil
}

// APIMergePullRequest implements PullRequestService.
func (p *pullRequestServiceImpl) APIMergePullRequest(pullRequestID string) error {
	// reuse internal MergePullRequest which returns domain error / nil
	return p.MergePullRequest(pullRequestID)
}

// APIReassignPullRequest implements PullRequestService.
func (p *pullRequestServiceImpl) APIReassignPullRequest(request models.ReassignRequest) (*models.ReassignResponse, error) {
	c := context.TODO()

	tx, err := p.pool.Begin(c)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(c)

	q := db.New(tx)
	qtx := q.WithTx(tx)

	pr, err := qtx.GetPullRequestByID(c, request.PullRequestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("pr not found")
		}
		return nil, err
	}

	if pr.Status == 1 {
		return nil, errors.New("cannot reassign on merged PR")
	}

	review, err := qtx.GetReviewByPRAndUser(c, db.GetReviewByPRAndUserParams{
		PullRequestID: request.PullRequestID,
		UserID:        request.OldReviewerID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("reviewer is not assigned to this PR")
		}
		return nil, err
	}

	oldUser, err := qtx.GetUserByID(c, request.OldReviewerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	teamUsers, err := qtx.GetUsersByTeamID(c, oldUser.TeamID)
	if err != nil {
		return nil, err
	}

	currentReviews, err := qtx.GetReviewsByPullRequestID(c, request.PullRequestID)
	if err != nil {
		return nil, err
	}
	assignedSet := make(map[string]struct{})
	for _, r := range currentReviews {
		assignedSet[r.UserID] = struct{}{}
	}

	var candidate *db.User
	for _, u := range teamUsers {
		if u.ID == oldUser.ID {
			continue
		}
		if !u.IsActive {
			continue
		}
		if u.ID == pr.AuthorID {
			continue
		}
		if _, taken := assignedSet[u.ID]; taken {
			continue
		}
		candidate = &u
		break
	}

	if candidate == nil {
		return nil, errors.New("no active replacement candidate in team")
	}

	if _, err := qtx.UpdateReviewUser(c, db.UpdateReviewUserParams{
		ID:     review.ID,
		UserID: candidate.ID,
	}); err != nil {
		return nil, err
	}

	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	latestPR, err := q.GetPullRequestByID(context.TODO(), request.PullRequestID)
	if err != nil {
		return nil, err
	}
	reviewsAfter, err := q.GetReviewsByPullRequestID(context.TODO(), request.PullRequestID)
	if err != nil {
		return nil, err
	}
	var assignedIDs []string
	for _, r := range reviewsAfter {
		assignedIDs = append(assignedIDs, r.UserID)
	}

	prDTO := dto.PullRequestDTO{
		ID:                latestPR.ID,
		PullRequestName:   latestPR.PullRequestName,
		AuthorID:          latestPR.AuthorID,
		Status:            enum.PrStatus(latestPR.Status),
		AssignedReviewers: assignedIDs,
	}

	return &models.ReassignResponse{
		PR:           prDTO,
		ReplacedByID: candidate.ID,
	}, nil
}

// CreatePullRequest implements PullRequestService.
func (p *pullRequestServiceImpl) CreatePullRequest(pullRequestID string, pullRequestName string, authorID string) (*db.PullRequest, error) {
	q := db.New(p.pool)
	ctx := context.TODO()

	if _, err := q.GetPullRequestByID(ctx, pullRequestID); err == nil {
		return nil, errors.New("PR id already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	author, err := q.GetUserByID(ctx, authorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("author not found")
		}
		return nil, err
	}

	users, err := q.GetUsersByTeamID(ctx, author.TeamID)
	if err != nil {
		return nil, err
	}
	var assigned []db.User
	for _, u := range users {
		if u.ID == author.ID {
			continue
		}
		if !u.IsActive {
			continue
		}
		assigned = append(assigned, u)
		if len(assigned) == 2 {
			break
		}
	}

	pr, err := q.CreatePullRequest(ctx, db.CreatePullRequestParams{
		ID:              pullRequestID,
		PullRequestName: pullRequestName,
		AuthorID:        authorID,
		Status:          0,
	})
	if err != nil {
		return nil, err
	}

	for _, a := range assigned {
		reviewID := uuid.New().String()
		if _, err := q.CreateReview(ctx, db.CreateReviewParams{
			ID:            reviewID,
			UserID:        a.ID,
			PullRequestID: pr.ID,
			Reviewed:      false,
		}); err != nil {
			return nil, err
		}
	}

	return &pr, nil
}

// MergePullRequest implements PullRequestService.
func (p *pullRequestServiceImpl) MergePullRequest(pullRequestID string) error {
	q := db.New(p.pool)
	ctx := context.TODO()

	pr, err := q.GetPullRequestByID(ctx, pullRequestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("pr not found")
		}
		return err
	}

	if pr.Status == 1 {
		return nil
	}

	mergedAt := time.Now().UTC()

	if _, err := q.UpdatePullRequestStatus(ctx, db.UpdatePullRequestStatusParams{
		ID:       pullRequestID,
		Status:   1,
		MergedAt: pgtype.Timestamptz{Time: mergedAt, Valid: true},
	}); err != nil {
		return err
	}

	return nil
}

func (p *pullRequestServiceImpl) ReassignPullRequest(pullRequestID string, oldReviewerID string) error {
	req := models.ReassignRequest{
		PullRequestID: pullRequestID,
		OldReviewerID: oldReviewerID,
	}
	_, err := p.APIReassignPullRequest(req)
	return err
}

func NewPullRequestService(pool *pgxpool.Pool) PullRequestService {
	return &pullRequestServiceImpl{
		pool: pool,
	}
}
