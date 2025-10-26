package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"mimic/backend/database"
)

func Validate(tokenString string) bool {
	tokenId := tokenString[7:]
	user, err := database.GetUserByToken(tokenId)
	return err == nil && user.ID != 0
}

func CreateToken(username string, password string) {
	token := GenerateSecureToken(32)
	userId, err := database.GetUserID(username)
	if err != nil {
		log.Fatal("Error getting user ID:", err)
	}
	err = database.UpdateSessionToken(userId, token)
	if err != nil {
		log.Fatal("Error updating session token:", err)
	}
	fmt.Printf("Creating token for %s %s", username, token)
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
