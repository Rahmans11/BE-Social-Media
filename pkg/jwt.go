package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTClaims(id int, role string) *JWTClaims {
	return &JWTClaims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func (jc *JWTClaims) GenToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("no secret found")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jc)
	return token.SignedString([]byte(jwtSecret))
}

func (jc *JWTClaims) VerifyToken(token string) (bool, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return false, errors.New("no secret found")
	}
	jwtToken, err := jwt.ParseWithClaims(token, jc, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return false, err
	}
	if !jwtToken.Valid {
		return false, jwt.ErrTokenExpired
	}
	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return false, err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return false, jwt.ErrTokenInvalidIssuer
	}
	return true, nil
}
