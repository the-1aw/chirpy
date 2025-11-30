package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/the-1aw/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func fromDbChirp(chirp database.Chirp) Chirp {
	return Chirp{
		ID:        chirp.ID,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
	}
}

func getProfaneWords() [3]string { return [...]string{"kerfuffle", "sharbert", "fornax"} }

func sanitizeChirpBody(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}
	cleanedWords := []string{}
	for word := range strings.SplitSeq(body, " ") {
		isProfaneWord := false
		for _, profaneWord := range getProfaneWords() {
			if profaneWord == strings.ToLower(word) {
				isProfaneWord = true
			}
		}
		if isProfaneWord {
			cleanedWords = append(cleanedWords, "****")
		} else {
			cleanedWords = append(cleanedWords, word)
		}
	}
	return strings.Join(cleanedWords, " "), nil
}

func getCreateChirpHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}
		type responseBody struct {
			CleanedBody string `json:"cleaned_body"`
		}
		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithErrorJSON(w, http.StatusBadRequest, err)
			return
		}

		body := requestBody{}
		if err := json.Unmarshal(rawBody, &body); err != nil {
			respondWithErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		chirpBody, err := sanitizeChirpBody(body.Body)
		if err != nil {
			respondWithErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			UserID: body.UserID,
			Body:   chirpBody,
		})
		if err != nil {
			respondWithErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		respondWithJSON(w, http.StatusCreated, fromDbChirp(chirp))
	})
}

func getGetChirpHandler(cfg *apiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type responseBody []Chirp
		chirpsFromDb, err := cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithErrorJSON(w, http.StatusInternalServerError, err)
			return
		}
		chirps := []Chirp{}
		for _, c := range chirpsFromDb {
			chirps = append(chirps, fromDbChirp(c))
		}
		respondWithJSON(w, http.StatusOK, chirps)
	})
}
