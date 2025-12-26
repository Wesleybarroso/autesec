package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (j *JWTManager) Generate(userID, certFingerprint string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"cert_fp": certFingerprint,
		"exp":      time.Now().Add(j.ttl).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) Validate(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, _ := claims["sub"].(string)
	certFP, _ := claims["cert_fp"].(string)

	if userID == "" || certFP == "" {
		return "", "", errors.New("invalid claims")
	}

	return userID, certFP, nil
}
