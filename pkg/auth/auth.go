package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	Checksum [32]byte
	Auth     bool
	jwt.RegisteredClaims
}

func passOK(input string) bool {
	secret := os.Getenv("TODO_PASSWORD")
	return subtle.ConstantTimeCompare(
		[]byte(input),
		[]byte(secret),
	) == 1
}

func buildJWT() (string, error) {
	claims := claims{
		Checksum: sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD") + "salt")),
		Auth:     true,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
}

func ValidateJWT(jwtData string) error {
	claims := claims{}
	parsedToken, err := jwt.ParseWithClaims(jwtData, &claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing mothod: %v", token.Method)
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})

	if err != nil || !parsedToken.Valid {
		return fmt.Errorf("error validating token")
	}

	if claims.Checksum != sha256.Sum256([]byte(os.Getenv("TODO_PASSWORD")+"salt")) {
		return fmt.Errorf("invalid hash")
	}

	return nil
}

func SignIn(password string) (string, error) {
	if !passOK(password) {
		return "", fmt.Errorf("invalid password")
	}

	return buildJWT()
}
