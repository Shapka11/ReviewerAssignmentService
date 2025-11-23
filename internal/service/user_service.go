package service

import (
	"context"
	"errors"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/repository"
)

var ErrUserFound = errors.New("user not found")

type userServiceImpl struct {
	userRepository repository.UserRepository
	prRepo         repository.PRRepository
}

func NewUserService(userRepository repository.UserRepository, prRepo repository.PRRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
		prRepo:         prRepo,
	}
}

func (s *userServiceImpl) SetIsActive(ctx context.Context, userID string, isActive bool) error {
	return s.userRepository.UpdateActivity(ctx, userID, isActive)
}

func (s *userServiceImpl) GetUserPRs(ctx context.Context, userID string) ([]*domains.PullRequestShort, error) {
	exists, err := s.userRepository.Exists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrUserFound
	}

	return s.prRepo.GetByReviewer(ctx, userID)
}

func (s *userServiceImpl) GetGlobalStats(ctx context.Context) (*domains.GlobalStats, error) {
	usersCount, err := s.userRepository.Count(ctx)
	if err != nil {
		return nil, err
	}

	prsCount, err := s.prRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &domains.GlobalStats{
		TotalUsers: usersCount,
		TotalPRs:   prsCount,
	}, nil
}
