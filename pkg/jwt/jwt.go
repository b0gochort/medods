package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"medods/model"
	"time"
)

func GenerateToken(ip, email string, expiration time.Duration) (string, error) {
	secretKey := viper.GetString("jwt.secret")
	claims := &model.JwtCustomClaims{
		IP:    ip,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
