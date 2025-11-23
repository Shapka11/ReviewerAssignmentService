package tests

import (
	"context"
	"testing"

	"ReviewerAssignmentService/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/service"
)

func TestService_CreatePR(t *testing.T) {
	type args struct {
		input domains.PullRequestInput
	}
	tests := []struct {
		name       string
		args       args
		setupMocks func(prRepo *mocks.PRRepository, userRepo *mocks.UserRepository)
		wantErr    bool
	}{
		{
			name: "Success",
			args: args{input: domains.PullRequestInput{ID: "pr-1", AuthorID: "u1"}},
			setupMocks: func(prRepo *mocks.PRRepository, userRepo *mocks.UserRepository) {
				userRepo.On("GetByID", mock.Anything, "u1").Return(&domains.User{ID: "u1", TeamName: "A"}, nil)
				userRepo.On("GetRandomActiveUsersByTeam", mock.Anything, "A", "u1", 2).Return([]string{"u2", "u3"}, nil)
				prRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mPR := mocks.NewPRRepository(t)
			mUser := mocks.NewUserRepository(t)
			tt.setupMocks(mPR, mUser)

			s := service.NewPRService(mPR, mUser)
			_, err := s.CreatePR(context.Background(), tt.args.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_MergePR(t *testing.T) {
	mPR := mocks.NewPRRepository(t)
	mUser := mocks.NewUserRepository(t)
	s := service.NewPRService(mPR, mUser)

	pr := &domains.PullRequest{ID: "pr-1", Status: domains.PRStatusOpen}
	mPR.On("GetByID", mock.Anything, "pr-1").Return(pr, nil)
	mPR.On("Update", mock.Anything, mock.MatchedBy(func(p *domains.PullRequest) bool {
		return p.Status == domains.PRStatusMerged
	})).Return(nil)

	got, err := s.MergePR(context.Background(), "pr-1")

	assert.NoError(t, err)
	assert.Equal(t, domains.PRStatusMerged, got.Status)
}
