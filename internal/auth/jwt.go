package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(playerID int64) (string, error) {
	claim := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(playerID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claim,
	)

	secret := []byte(os.Getenv("JWT_SECRET"))

	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
	)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	playerID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid player id in token")
	}

	return playerID, nil
}
