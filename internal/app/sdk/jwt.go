package sdk

import (
	"fmt"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store/types"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	AccessTokenKey []byte
	AccessTokenTTL time.Duration
}

type JWTClaims struct {
	UserId  int64 `json:"uid"`
	IsAdmin bool  `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		AccessTokenKey: []byte(secret),
		AccessTokenTTL: ttl,
	}
}

func (j *JWTManager) GenerateToken(user *types.User) (string, error) {
	claims := &JWTClaims{
		UserId:  user.ID,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.AccessTokenTTL)),
			Issuer:    "critiquefi-service",
			Subject:   fmt.Sprint(user.ID),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.AccessTokenKey)
}

func (j *JWTManager) ParseToken(token string) (*JWTClaims, error) {
	t, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.AccessTokenKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*JWTClaims)
	if !ok || !t.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
