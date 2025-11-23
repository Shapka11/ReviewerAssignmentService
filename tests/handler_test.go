package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/handler"
	"ReviewerAssignmentService/mocks"
)

func TestHandler_CreatePR(t *testing.T) {
	tests := []struct {
		name   string
		body   map[string]interface{}
		mock   func(s *mocks.PRService)
		status int
	}{
		{
			name: "Success 201",
			body: map[string]interface{}{
				"pull_request_id":   "pr-1",
				"pull_request_name": "Test",
				"author_id":         "u1",
			},
			mock: func(s *mocks.PRService) {
				s.On("CreatePR", mock.Anything, mock.Anything).
					Return(&domains.PullRequest{ID: "pr-1"}, nil)
			},
			status: http.StatusCreated,
		},
		{
			name:   "Bad Request",
			body:   nil,
			mock:   func(s *mocks.PRService) {},
			status: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			body: map[string]interface{}{
				"pull_request_id": "pr-1", "pull_request_name": "T", "author_id": "u1",
			},
			mock: func(s *mocks.PRService) {
				s.On("CreatePR", mock.Anything, mock.Anything).
					Return(nil, errors.New("db error"))
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prService := mocks.NewPRService(t)
			tt.mock(prService)

			h := handler.New(mocks.NewTeamService(t), mocks.NewUserService(t), prService)
			router := h.InitRoutes()

			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
			router.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
