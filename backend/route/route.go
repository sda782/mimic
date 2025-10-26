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
)

const maxUploadSize = 10 << 30 // 10 GB

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading multipart data: %v", err), http.StatusBadRequest)
		return
	}

	userID, filePart, filename, err := parseFormFieldsAndFile(mr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userID == -1 {
		http.Error(w, "User not found", http.StatusNotFound)
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
		UserID:    userID,
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
		"status":  "success",
		"message": fmt.Sprintf("Upload complete, download your file here: %s://%s/%s", scheme, r.Host, fileUpload.ShortCode),
		"url":     fmt.Sprintf("%s://%s/%s", scheme, r.Host, fileUpload.ShortCode),
	}

	json.NewEncoder(w).Encode(response)
}

func GetUploads(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	token = token[7:]
	user, err := database.GetUserByToken(token)
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

func parseFormFieldsAndFile(mr *multipart.Reader) (userID int, filePart multipart.Part, filename string, err error) {
	userID = -1
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
			if part.FormName() == "name" {
				value, _ := io.ReadAll(part)
				userID, e = database.GetUserID(string(value))
				if e != nil {
					err = fmt.Errorf("error getting user ID: %v", e)
					return
				}
			}
			continue
		}

		filename = part.FileName()
		filePart = *part
		break
	}
	return
}
