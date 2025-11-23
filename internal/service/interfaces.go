package service

import (
	"context"

	"ReviewerAssignmentService/internal/domains"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *domains.Team) (*domains.Team, error)
	GetTeam(ctx context.Context, name string) (*domains.Team, error)
}

type UserService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) error
	GetUserPRs(ctx context.Context, userID string) ([]*domains.PullRequestShort, error)
	GetGlobalStats(ctx context.Context) (*domains.GlobalStats, error)
}

type PRService interface {
	CreatePR(ctx context.Context, input domains.PullRequestInput) (*domains.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*domains.PullRequest, error)
	UpdateReviewer(ctx context.Context, prID string, oldReviewerID string) (*domains.PullRequest, string, error)
}
