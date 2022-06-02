package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

// JWTTokenVerifier verifies token
type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify verifies token and returns accountID
func (v *JWTTokenVerifier) Verify(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("cannot parse token: %v", err)
	}

	if !t.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("cannot get claims")
	}

	// token is expired
	if err := claims.Valid(); err != nil {
		return "", fmt.Errorf("claims is invalid: %v", err)
	}

	return claims.Subject, nil
}
