package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"toll-booth/internal/db"
	"toll-booth/internal/handler"
	"toll-booth/internal/middleware"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handleRoot)
	mux.HandleFunc("/register", handler.HandleRegisterUser)
	mux.HandleFunc("/login", handler.HandleLoginUser)
	mux.HandleFunc("/reset-password", handler.HandleResetPassword)
	mux.HandleFunc("/shortner", handler.HandleLinkShortner)
	mux.HandleFunc("/redirect", handler.HandleRedirect)
	mux.HandleFunc("/redirect-request/{shortCode}", handler.HandleRedirectRequest)
	mux.HandleFunc("/getLinks", handler.GetAllLinks)

	srv := &http.Server{
		Addr:         port,
		Handler:      middleware.CORS(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on %s", port)
	log.Fatal(srv.ListenAndServe())

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World!"))

	if err != nil {
		slog.Error("Couldn't write root", "err", err)
	}
}
