package model

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtCustomClaims struct {
	IP           string `json:"ip"`
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
	jwt.RegisteredClaims
}

func (c *JwtCustomClaims) Valid() error {
	expirationTime := c.ExpiresAt.Time.Unix()

	if time.Now().Unix() > expirationTime {
		return errors.New("token has expired")
	}

	return nil
}
