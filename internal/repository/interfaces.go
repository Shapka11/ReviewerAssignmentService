package repository

import (
	"context"

	"ReviewerAssignmentService/internal/domains"
)

type UserRepository interface {
	Create(ctx context.Context, user *domains.User) error
	Exists(ctx context.Context, userName string) (bool, error)
	GetByID(ctx context.Context, id string) (*domains.User, error)
	GetRandomActiveUsersByTeam(ctx context.Context, teamName string, excludeUserID string, limit int) ([]string, error)
	UpdateActivity(ctx context.Context, userID string, isActive bool) error
	DeactivateTeamMembers(ctx context.Context, teamName string) error
	Count(ctx context.Context) (int, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *domains.Team) error
	Exists(ctx context.Context, teamName string) (bool, error)
	GetByName(ctx context.Context, teamName string) (*domains.Team, error)
}

type PRRepository interface {
	Create(ctx context.Context, pr *domains.PullRequest) error
	Exists(ctx context.Context, prName string) (bool, error)
	GetByID(ctx context.Context, id string) (*domains.PullRequest, error)
	GetByReviewer(ctx context.Context, reviewerID string) ([]*domains.PullRequestShort, error)
	Update(ctx context.Context, pr *domains.PullRequest) error
	Count(ctx context.Context) (int, error)
}
