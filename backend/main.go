package main

import (
	"fmt"
	"mimic/backend/auth"
	"mimic/backend/database"
	"mimic/backend/misc"
	"mimic/backend/route"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	load()
	setupAdminUser()

	mux := http.NewServeMux()

	mux.HandleFunc("/uploads", auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			route.GetUploads(w, r)
		case "POST":
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/upload", auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			route.HandleUpload(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			route.Login(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/session/validate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			route.ValidateSession(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			route.HandleShortCode(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	handler := auth.CorsMiddleware(mux)

	fmt.Println("Server running on port", port)
	http.ListenAndServe(":"+port, handler)

	database.Close()
}
func load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	misc.Init()
	database.Init()
}

func setupAdminUser() {
	pass := os.Getenv("ADMIN_PASSWORD")
	if pass == "" {
		fmt.Println("Admin password not set, skipping setup")
		return
	}
	_, err := database.GetUser("admin", pass)
	if err != nil {
		_, err = database.CreateUser("admin", pass)
		if err != nil {
			fmt.Println("Error creating admin user:", err)
			return
		}
	}
}
