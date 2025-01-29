package util

import (
	"fmt"
	"os"
	modelsuser "park/models/modelsUser"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateJWT(userID int, username string, role modelsuser.RoleType, parkno string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	secretKey := os.Getenv("SECRET_KEY_JWT")
	claims := jwt.MapClaims{
		"user_id":  strconv.Itoa(userID),
		"username": username,
		"role":     role,
		"parkno":   parkno,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
