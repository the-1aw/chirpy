package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/the-1aw/chirpy/internal/auth"
)

type PolkaEventData struct {
	UserID uuid.UUID `json:"user_id"`
}

type PolkaEvent struct {
	Event string         `json:"event"`
	Data  PolkaEventData `json:"data"`
}

const (
	EventUserUpgrade = "user.upgraded"
)

func getPolkaWebhookHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiKey, err := auth.GetApiKey(r.Header); apiKey != cfg.polkaKey {
			respondWithErrorJSON(w, http.StatusUnauthorized, fmt.Errorf("Unauthorized %v", err))
			return
		}
		decoder := json.NewDecoder(r.Body)
		event := PolkaEvent{}
		if err := decoder.Decode(&event); err != nil {
			respondWithErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		if event.Event != EventUserUpgrade {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if err := cfg.db.UpgradeUserToChirpyRed(r.Context(), event.Data.UserID); err != nil {
			respondWithErrorJSON(w, http.StatusNotFound, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
