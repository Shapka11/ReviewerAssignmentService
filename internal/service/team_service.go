package service

import (
	"context"
	"errors"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/repository"
)

var ErrTeamExists = errors.New("team already exists")

type teamServiceImpl struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamServiceImpl{teamRepository: repo}
}

func (s *teamServiceImpl) CreateTeam(ctx context.Context, team *domains.Team) (*domains.Team, error) {
	exists, err := s.teamRepository.Exists(ctx, team.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTeamExists
	}

	if err := s.teamRepository.Create(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *teamServiceImpl) GetTeam(ctx context.Context, name string) (*domains.Team, error) {
	return s.teamRepository.GetByName(ctx, name)
}
