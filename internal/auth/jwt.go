package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenData struct {
	PlayerID  int64
	TokenID   string
	ExpiresAt time.Time
}

func GenerateToken(playerID int64) (string, error) {
	claim := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(playerID, 10),
		ID:        uuid.NewString(),
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

func ValidateToken(tokenString string) (TokenData, error) {
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
		return TokenData{}, err
	}

	if !token.Valid {
		return TokenData{}, fmt.Errorf("invalid token")
	}

	playerID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return TokenData{}, fmt.Errorf("invalid player id in token")
	}
	
	if claims.ID == "" {
		return TokenData{}, fmt.Errorf("token id is missing")
	}

	if claims.ExpiresAt == nil {
		return TokenData{}, fmt.Errorf("token expiration is missing")
	}

	return TokenData{
		PlayerID:  playerID,
		TokenID:   claims.ID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
