package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"ReviewerAssignmentService/internal/domains"
	"ReviewerAssignmentService/internal/service"
)

const ErrDuplicateKeyValue = "duplicate key value"

type createPRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type prIDRequest struct {
	ID string `json:"pull_request_id"`
}

type reassignRequest struct {
	ID        string `json:"pull_request_id"`
	OldUserID string `json:"old_user_id"`
}

func (h *Handler) createPR(w http.ResponseWriter, r *http.Request) {
	var req createPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgInvalidJSON)
		return
	}

	input := domains.PullRequestInput{
		ID:       req.ID,
		Name:     req.Name,
		AuthorID: req.AuthorID,
	}

	pr, err := h.prService.CreatePR(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrAuthorNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, ErrMsgAuthorNotFound)
			return
		}

		if strings.Contains(err.Error(), ErrDuplicateKeyValue) {
			writeError(w, http.StatusConflict, ErrCodePRExists, ErrMsgPRExists)
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"pr": pr,
	})
}

func (h *Handler) mergePR(w http.ResponseWriter, r *http.Request) {
	var req prIDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgInvalidJSON)
		return
	}

	pr, err := h.prService.MergePR(r.Context(), req.ID)
	if err != nil {
		if errors.Is(err, service.ErrPRNotFound) {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, ErrMsgPRNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pr": pr,
	})
}

func (h *Handler) reassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req reassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgInvalidJSON)
		return
	}

	pr, newID, err := h.prService.UpdateReviewer(r.Context(), req.ID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPRNotFound):
			writeError(w, http.StatusNotFound, ErrCodeNotFound, ErrMsgPRNotFound)
		case errors.Is(err, service.ErrPRMerged):
			writeError(w, http.StatusConflict, ErrCodePRMerged, ErrMsgPRMerged)
		case errors.Is(err, service.ErrReviewerNotAssigned):
			writeError(w, http.StatusConflict, ErrCodeNotAssigned, ErrMsgReviewerNotAssigned)
		case errors.Is(err, service.ErrNoCandidates):
			writeError(w, http.StatusConflict, ErrCodeNoCandidate, ErrMsgNoCandidate)
		default:
			writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": newID,
	})
}
