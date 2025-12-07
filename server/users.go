package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/the-1aw/chirpy/internal/auth"
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

type UserContent struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getCreateUserHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		req := UserContent{}
		if err := decoder.Decode(&req); err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		if len(req.Password) == 0 {
			respondWithErrorJSON(w, http.StatusBadRequest, errors.New("Password is required"))
			return
		}
		hashed_password, err := auth.HashPassword(req.Password)
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
			Email:          req.Email,
			HashedPassword: hashed_password,
		})
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusCreated, fromDbUser(user))
	})
}

func getUpdateUserHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		uid, err := auth.ValidateJWTFromHeader(r.Header, cfg.jwtSecret)
		if err != nil {
			respondWithErrorJSON(w, http.StatusUnauthorized, err)
		}
		body := UserContent{}
		if err := decoder.Decode(&body); err != nil || body.Email == "" || body.Password == "" {
			respondWithErrorJSON(w, http.StatusBadRequest, fmt.Errorf("Bad request %v\n", err))
			return
		}
		hashed_password, err := auth.HashPassword(body.Password)
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		err = cfg.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
			ID:             uid,
			Email:          body.Email,
			HashedPassword: hashed_password,
		})
		user, err := cfg.db.GetUserByEmail(r.Context(), body.Email)
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, fromDbUser(user))
	})
}
