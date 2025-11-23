package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/the-1aw/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func fromDbUser(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

func getCreateUserHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			Email string `json:"email"`
		}
		decoder := json.NewDecoder(r.Body)
		req := requestBody{}
		if err := decoder.Decode(&req); err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		user, err := cfg.db.CreateUser(r.Context(), req.Email)
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusCreated, fromDbUser(user))
	})
}
