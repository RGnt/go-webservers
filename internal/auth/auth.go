package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	stringHash, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return "", err
	}

	return string(stringHash), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJSWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	type customClaims struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	userId, err := uuid.Parse(subject)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("auth info missing from the header")
	}

	return strings.Replace(bearerToken, "Bearer ", "", -1), nil
}

func MakeRefreshToken() (string, error) {
	randSlice := make([]byte, 32)
	_, err := rand.Read(randSlice)
	if err != nil {
		fmt.Println("error, ", err)
		return "", err
	}

	return hex.EncodeToString(randSlice), nil
}
