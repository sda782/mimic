package fs

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func SaveTempFile(part multipart.Part) (string, error) {
	tempFile, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return "", err
	}
	tempPath := tempFile.Name()
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, &part); err != nil {
		os.Remove(tempPath)
		return "", err
	}

	return tempPath, nil
}

func MoveFileToFinal(tempPath, filename string) (string, error) {
	finalPath := fmt.Sprintf("%s%s", os.Getenv("UPLOAD_PATH"), filename)

	src, err := os.Open(tempPath)
	if err != nil {
		return "", err
	}

	dst, err := os.Create(finalPath)
	if err != nil {
		src.Close()
		return "", err
	}

	if _, err := io.Copy(dst, src); err != nil {
		src.Close()
		dst.Close()
		return "", err
	}

	src.Close()
	dst.Close()

	if err := os.Remove(tempPath); err != nil {
		return "", err
	}

	return finalPath, nil
}
