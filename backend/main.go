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

	http.HandleFunc("/uploads/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if auth.Validate(r.Header.Get("Authorization")) {
				route.GetUploads(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		case "POST":
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(os.Getenv("WEBAPP_PATH")))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			route.HandleShortCode(w, r)
		case "POST":
			if auth.Validate(r.Header.Get("Authorization")) {
				route.HandleUpload(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
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
