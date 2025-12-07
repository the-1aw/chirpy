package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/the-1aw/chirpy/internal/auth"
	"github.com/the-1aw/chirpy/internal/database"
)

func getLoginHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		type responseBody struct {
			User
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}
		decoder := json.NewDecoder(r.Body)
		req := requestBody{}
		if err := decoder.Decode(&req); err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			respondWithErrorJSON(w, http.StatusUnauthorized, err)
			return
		}
		passMatches, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
		if err != nil || !passMatches {
			respondWithErrorJSON(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
		rt, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     auth.MakeRefreshToken(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
			UserID:    user.ID,
		})
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, responseBody{
			User:         fromDbUser(user),
			Token:        token,
			RefreshToken: rt.Token,
		})
	})
}

func getRefreshHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithErrorJSON(w, http.StatusUnauthorized, err)
			return
		}
		rt, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
		if err != nil {
			respondWithErrorJSON(w, http.StatusUnauthorized, err)
			return
		}
		token, err := auth.MakeJWT(rt.UserID, cfg.jwtSecret, time.Hour)
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		respondWithJSON(w, http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
	})
}

func getRevokeHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithErrorJSON(w, http.StatusUnauthorized, err)
			return
		}
		err = cfg.db.RevokeToken(r.Context(), refreshToken)
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
