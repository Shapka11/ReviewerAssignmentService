package tests

import (
	"context"
	"testing"
	"time"

	"ReviewerAssignmentService/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/service"
)

func TestPRService_CreatePR(t *testing.T) {
	testCases := []struct {
		name              string
		input             domains.PullRequestInput
		setupMocks        func(pr *mocks.PRRepository, u *mocks.UserRepository)
		expectError       bool
		expectedReviewers []string
	}{
		{
			name: "OK: 2 revs",
			input: domains.PullRequestInput{
				ID:       "pr1",
				Name:     "Feature 1",
				AuthorID: "u1",
			},
			setupMocks: func(pr *mocks.PRRepository, u *mocks.UserRepository) {
				u.On("GetByID", mock.Anything, "u1").Return(&domains.User{
					ID:       "u1",
					TeamName: "T1",
				}, nil)

				u.On("GetRandomActiveUsersByTeam", mock.Anything, "T1", "u1", 2).
					Return([]string{"r1", "r2"}, nil)

				pr.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError:       false,
			expectedReviewers: []string{"r1", "r2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockPRRepo := mocks.NewPRRepository(t)
			mockUserRepo := mocks.NewUserRepository(t)
			tc.setupMocks(mockPRRepo, mockUserRepo)

			res, err := service.NewPRService(mockPRRepo, mockUserRepo).CreatePR(context.Background(), tc.input)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expectedReviewers, res.AssignedReviewers)
			}
		})
	}
}

func TestPRService_MergePR(t *testing.T) {
	testCases := []struct {
		name           string
		prID           string
		setup          func(prRepo *mocks.PRRepository)
		expectedStatus domains.PRStatus
		expectError    bool
	}{
		{
			name: "OK: Merge open PR",
			prID: "pr-open",
			setup: func(prRepo *mocks.PRRepository) {
				prRepo.On("GetByID", mock.Anything, "pr-open").Return(&domains.PullRequest{
					ID:     "pr-open",
					Status: domains.PRStatusOpen,
				}, nil)
				prRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: domains.PRStatusMerged,
			expectError:    false,
		},
		{
			name: "Idempotency: PR already merged",
			prID: "pr-merged",
			setup: func(prRepo *mocks.PRRepository) {
				now := time.Now()
				prRepo.On("GetByID", mock.Anything, "pr-merged").Return(&domains.PullRequest{
					ID:       "pr-merged",
					Status:   domains.PRStatusMerged,
					MergedAt: &now,
				}, nil)
			},
			expectedStatus: domains.PRStatusMerged,
			expectError:    false,
		},
	}

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			prRepo := mocks.NewPRRepository(t)
			ts.setup(prRepo)

			svc := service.NewPRService(prRepo, mocks.NewUserRepository(t))
			result, err := svc.MergePR(context.Background(), ts.prID)

			if ts.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, ts.expectedStatus, result.Status)
			}
		})
	}
}

func TestPRService_UpdateReviewer(t *testing.T) {
	testCases := []struct {
		name        string
		prID        string
		oldID       string
		setup       func(prRepo *mocks.PRRepository, userRepo *mocks.UserRepository)
		expectedID  string
		expectError error
	}{
		{
			name:  "OK: Reassign reviewer",
			prID:  "pr1",
			oldID: "old",
			setup: func(prRepo *mocks.PRRepository, userRepo *mocks.UserRepository) {
				prRepo.On("GetByID", mock.Anything, "pr1").Return(&domains.PullRequest{
					ID:                "pr1",
					AuthorID:          "author",
					Status:            domains.PRStatusOpen,
					AssignedReviewers: []string{"old", "second"},
				}, nil)

				userRepo.On("GetByID", mock.Anything, "old").Return(&domains.User{
					ID:       "old",
					TeamName: "Team",
				}, nil)

				userRepo.On("GetRandomActiveUsersByTeam", mock.Anything, "Team", "old", 5).
					Return([]string{"author", "second", "new"}, nil)

				prRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedID:  "new",
			expectError: nil,
		},
		{
			name:  "Fail: Cannot reassign on merged pr",
			prID:  "pr2",
			oldID: "u1",
			setup: func(prRepo *mocks.PRRepository, userRepo *mocks.UserRepository) {
				prRepo.On("GetByID", mock.Anything, "pr2").Return(&domains.PullRequest{
					ID:     "pr2",
					Status: domains.PRStatusMerged,
				}, nil)
			},
			expectedID:  "",
			expectError: service.ErrPRMerged,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prRepo := mocks.NewPRRepository(t)
			userRepo := mocks.NewUserRepository(t)
			tc.setup(prRepo, userRepo)

			svc := service.NewPRService(prRepo, userRepo)
			_, newID, err := svc.UpdateReviewer(context.Background(), tc.prID, tc.oldID)

			if tc.expectError != nil {
				assert.ErrorIs(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, newID)
			}
		})
	}
}
