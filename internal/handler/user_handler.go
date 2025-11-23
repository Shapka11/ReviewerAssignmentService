package handler

import (
	"encoding/json"
	"net/http"
)

type setActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) setUserActive(w http.ResponseWriter, r *http.Request) {
	var req setActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgInvalidJSON)
		return
	}

	err := h.userService.SetIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if err.Error() == ErrMsgUserNotFound {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, ErrMsgUserNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": map[string]interface{}{
			"user_id":   req.UserID,
			"is_active": req.IsActive,
		},
	})
}

func (h *Handler) getUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgMissingUserID)
		return
	}

	prs, err := h.userService.GetUserPRs(r.Context(), userID)
	if err != nil {
		if err.Error() == ErrMsgUserNotFound {
			writeError(w, http.StatusNotFound, ErrCodeNotFound, ErrMsgUserNotFound)
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.userService.GetGlobalStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
