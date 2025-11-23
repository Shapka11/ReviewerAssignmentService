package handler

import (
	"encoding/json"
	"net/http"

	"ReviewerAssignmentService/internal/domains"
)

func (h *Handler) createTeam(w http.ResponseWriter, r *http.Request) {
	var team domains.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgInvalidJSON)
		return
	}

	createdTeam, err := h.teamService.CreateTeam(r.Context(), &team)
	if err != nil {
		if err.Error() == ErrMsgTeamExists {
			writeError(w, http.StatusBadRequest, ErrCodeTeamExists, ErrMsgTeamExists)
			return
		}
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"team": createdTeam,
	})
}

func (h *Handler) getTeam(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("team_name")
	if name == "" {
		writeError(w, http.StatusBadRequest, ErrCodeBadRequest, ErrMsgMissingTeamName)
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, ErrCodeInternalError, err.Error())
		return
	}
	if team == nil {
		writeError(w, http.StatusNotFound, ErrCodeNotFound, ErrMsgTeamNotFound)
		return
	}

	writeJSON(w, http.StatusOK, team)
}
