package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/the-1aw/chirpy/server"
)

func main() {
	godotenv.Load()
	log.Fatal(server.Run())
}
