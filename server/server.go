package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/the-1aw/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverHits.Add(1)
			next.ServeHTTP(w, r)
		},
	)
}

func (cfg *apiConfig) requestCount(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PLATFORM") != "dev" {
		respondWithErrorJSON(w, http.StatusForbidden, fmt.Errorf("Cannot reset database if not in dev mode"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(r.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func Run() error {
	dbUrl := os.Getenv("DB_URL")
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return fmt.Errorf("Missing jwt secret")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	dbQueries := database.New(db)
	cfg := apiConfig{
		db:        dbQueries,
		jwtSecret: jwtSecret,
	}
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("POST /api/users", getCreateUserHandler(&cfg))
	mux.Handle("POST /api/login", getLoginHandler(&cfg))
	mux.Handle("POST /api/chirps", getCreateChirpHandler(&cfg))
	mux.Handle("GET /api/chirps", getGetChirpsHandler(&cfg))
	mux.Handle("GET /api/chirps/{chirpID}", getGetChirpByIdHandler(&cfg))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", cfg.requestCount)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	return server.ListenAndServe()
}
