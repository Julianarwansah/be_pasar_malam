package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
	expiry int
}

type Claims struct {
	UserID uint   `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, expiryHours int) *JWTService {
	return &JWTService{secret: []byte(secret), expiry: expiryHours}
}

func (s *JWTService) Generate(userID uint, email, role string) (string, int, error) {
	exp := time.Now().Add(time.Duration(s.expiry) * time.Hour)
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := tok.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}
	return str, s.expiry * 3600, nil
}

func (s *JWTService) Parse(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
