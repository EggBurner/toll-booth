package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"toll-booth/internal/db"
	"toll-booth/internal/handler"

	"github.com/joho/godotenv"
)

func main() {

	err := db.DBConnect()

	if err != nil {
		log.Fatalf("\nError while connecting to the database %v", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env vars")
	}

	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/register", handler.HandleRegisterUser)
	mux.HandleFunc("/login", handler.HandleLoginUser)
	mux.HandleFunc("/reset-password", handler.HandleResetPassword)
	mux.HandleFunc("/shortner", handler.HandleLinkShortner)
	mux.HandleFunc("/redirect", handler.HandleRedirect)
	mux.HandleFunc("/redirect-request/{shortCode}", handler.HandleRedirectRequest)

	log.Fatal(http.ListenAndServe(PORT, mux))

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World!"))

	if err != nil {
		slog.Error("Couldn't write root", "err", err)
	}
}
