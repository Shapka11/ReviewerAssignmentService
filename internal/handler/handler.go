package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"ReviewerAssignmentService/internal/service"
)

type Handler struct {
	teamService service.TeamService
	userService service.UserService
	prService   service.PRService
}

func New(team service.TeamService, user service.UserService, pr service.PRService) *Handler {
	return &Handler{
		teamService: team,
		userService: user,
		prService:   pr,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /team/add", h.createTeam)
	mux.HandleFunc("GET /team/get", h.getTeam)

	mux.HandleFunc("POST /users/setIsActive", h.setUserActive)
	mux.HandleFunc("GET /users/getReview", h.getUserReviews)

	mux.HandleFunc("POST /pullRequest/create", h.createPR)
	mux.HandleFunc("POST /pullRequest/merge", h.mergePR)
	mux.HandleFunc("POST /pullRequest/reassign", h.reassignReviewer)

	mux.HandleFunc("GET /stats", h.getStats)

	return mux
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	resp := errorResponse{
		Error: errorDetail{
			Code:    code,
			Message: message,
		},
	}
	writeJSON(w, status, resp)
}
