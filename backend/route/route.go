package route

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"mimic/backend/database"
	"mimic/backend/fs"
	"mimic/backend/types"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const maxUploadSize = 10 << 30 // 10 GB

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	cookie, err := r.Cookie("session_token")
	user, err := database.GetUserByToken(cookie.Value)
	fmt.Println(user.ID)
	if user.ID == -1 || err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading multipart data: %v", err), http.StatusBadRequest)
		return
	}

	filePart, filename, err := parseFormFieldsAndFile(mr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tempPath, err := fs.SaveTempFile(filePart)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving temp file: %v", err), http.StatusInternalServerError)
		return
	}

	defer func(path string) {
		if _, err := os.Stat(path); err == nil {
			os.Remove(path)
		}
	}(tempPath)

	finalPath, err := fs.MoveFileToFinal(tempPath, filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving final file: %v", err), http.StatusInternalServerError)
		return
	}

	fileUpload := types.Upload{
		UserID:    user.ID,
		FileName:  filename,
		FilePath:  finalPath,
		ShortCode: "",
	}

	err = database.Insert(&fileUpload)

	if err != nil {
		os.Remove(finalPath)
		http.Error(w, fmt.Sprintf("Database insert failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	response := map[string]interface{}{
		"status": "success",
		"url":    fmt.Sprintf("%s://%s/%s", scheme, r.Host, fileUpload.ShortCode),
	}

	json.NewEncoder(w).Encode(response)
}

func GetUploads(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	user, err := database.GetUserByToken(cookie.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user ID: %v", err), http.StatusUnauthorized)
		return
	}
	uploads := database.GetUploads(user.ID)

	mapped_uploads := make([]map[string]any, len(uploads))
	for i, upload := range uploads {
		mapped_uploads[i] = map[string]any{
			"short_code": upload.ShortCode,
			"filename":   upload.FileName,
		}
	}

	response := map[string]interface{}{
		"status": "success",
		"data":   mapped_uploads,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func HandleShortCode(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]
	if shortCode == "" || len(shortCode) != 4 {
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}
	upload := database.GetUpload(shortCode)

	http.ServeFile(w, r, upload.FilePath)
}

func parseFormFieldsAndFile(mr *multipart.Reader) (filePart multipart.Part, filename string, err error) {
	for {
		part, e := mr.NextPart()
		if e == io.EOF {
			break
		}
		if e != nil {
			err = fmt.Errorf("error reading multipart: %v", e)
			return
		}

		if part.FileName() == "" {
			continue
		}

		filename = part.FileName()
		filePart = *part
		break
	}
	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := database.GetUser(username, password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()

	err = database.UpdateSessionToken(user.ID, sessionToken)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Login successful",
	})
}

func ValidateSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No session", http.StatusUnauthorized)
		return
	}

	_, err = database.GetUserByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}
