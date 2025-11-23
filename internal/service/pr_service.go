package service

import (
	"context"
	"errors"
	"time"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/repository"
)

var (
	ErrPRNotFound               = errors.New("pull request not found")
	ErrPRMerged                 = errors.New("cannot edit merged PR")
	ErrReviewerNotAssigned      = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidates             = errors.New("no available candidates for assignment")
	ErrAuthorNotFound           = errors.New("author not found")
	ErrOriginalReviewerNotFound = errors.New("original reviewer user not found")
)

type prServiceImpl struct {
	prRepository   repository.PRRepository
	userRepository repository.UserRepository
}

func NewPRService(prRepository repository.PRRepository, userRepository repository.UserRepository) PRService {
	return &prServiceImpl{
		prRepository:   prRepository,
		userRepository: userRepository,
	}
}

func (s *prServiceImpl) CreatePR(ctx context.Context, input domains.PullRequestInput) (*domains.PullRequest, error) {
	author, err := s.userRepository.GetByID(ctx, input.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, ErrAuthorNotFound
	}

	candidateIDs, err := s.userRepository.GetRandomActiveUsersByTeam(ctx, author.TeamName, author.ID, 2)
	if err != nil {
		return nil, err
	}

	pr := &domains.PullRequest{
		ID:                input.ID,
		Name:              input.Name,
		AuthorID:          input.AuthorID,
		Status:            domains.PRStatusOpen,
		AssignedReviewers: candidateIDs,
	}

	if err := s.prRepository.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prServiceImpl) MergePR(ctx context.Context, prID string) (*domains.PullRequest, error) {
	pr, err := s.prRepository.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, ErrPRNotFound
	}

	if pr.Status == domains.PRStatusMerged {
		return pr, nil
	}

	pr.Status = domains.PRStatusMerged
	if err := s.prRepository.Update(ctx, pr); err != nil {
		return nil, err
	}

	now := time.Now()
	pr.MergedAt = &now

	return pr, nil
}

func (s *prServiceImpl) UpdateReviewer(
	ctx context.Context, prID string, oldReviewerID string) (*domains.PullRequest, string, error) {
	pr, err := s.prRepository.GetByID(ctx, prID)
	if err != nil {
		return nil, "", err
	}
	if pr == nil {
		return nil, "", ErrPRNotFound
	}

	if pr.Status == domains.PRStatusMerged {
		return nil, "", ErrPRMerged
	}

	idx := -1
	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, "", ErrReviewerNotAssigned
	}

	oldUser, err := s.userRepository.GetByID(ctx, oldReviewerID)
	if err != nil {
		return nil, "", err
	}
	if oldUser == nil {
		return nil, "", ErrOriginalReviewerNotFound
	}

	candidates, err := s.userRepository.GetRandomActiveUsersByTeam(ctx, oldUser.TeamName, oldReviewerID, 5)
	if err != nil {
		return nil, "", err
	}

	var newReviewerID string
	for _, candID := range candidates {
		if candID == pr.AuthorID {
			continue
		}

		alreadyAssigned := false
		for _, assigned := range pr.AssignedReviewers {
			if assigned == candID {
				alreadyAssigned = true
				break
			}
		}
		if !alreadyAssigned {
			newReviewerID = candID
			break
		}
	}

	if newReviewerID == "" {
		return nil, "", ErrNoCandidates
	}

	pr.AssignedReviewers[idx] = newReviewerID

	if err := s.prRepository.Update(ctx, pr); err != nil {
		return nil, "", err
	}

	return pr, newReviewerID, nil
}
