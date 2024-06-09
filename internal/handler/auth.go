package handler

import (
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log/slog"
	"medods/internal/repository/postgres"
	"medods/model"
	"medods/pkg/jwt"
	"net/http"
	"net/smtp"
	"time"
)

const (
	refreshTokenLife = 24 * time.Hour
	accessTokenLife  = 4 * time.Hour
)

func Login(log *slog.Logger, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var guid string
		guid = r.URL.Query().Get("guid")

		userIP := r.Header.Get("X-Forwarded-For")

		email, err := postgres.GetEmail(guid, db)

		if guid == "" || userIP == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		accessToken, err := jwt.GenerateToken(userIP, email, accessTokenLife)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		refreshToken, err := jwt.GenerateToken(userIP, email, refreshTokenLife)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cookie := http.Cookie{Name: "access_token",
			Value:   accessToken,
			Expires: time.Now().Add(refreshTokenLife)}
		http.SetCookie(w, &cookie)

		cookie = http.Cookie{Name: "refresh_token",
			Value:   refreshToken,
			Expires: time.Now().Add(refreshTokenLife)}
		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secretKey := viper.GetString("jwt.secret_key")

		userIP := r.Header.Get("X-Forwarded-For")

		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			slog.Info("refresh token cookie error", err)
			return
		}

		refreshToken := refreshTokenCookie.Value

		claims := model.JwtCustomClaims{}

		if claims.IP != userIP {
			if err := sendMail(claims.Email); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				slog.Info("refresh token send mail error", err)
				return
			}

			return
		}

		token, err := jwt2.ParseWithClaims(refreshToken,
			&claims,
			func(token *jwt2.Token) (interface{}, error) {
				return secretKey, nil
			})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			slog.Info("refresh token parse error", err)
			return
		}

		newAccessToken, err := jwt.GenerateToken(userIP, claims.Email, accessTokenLife)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Info("refresh token generate error", err)
			return
		}
		newRefreshToken, err := jwt.GenerateToken(userIP, claims.Email, refreshTokenLife)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			slog.Info("refresh token generate error", err)
			return
		}

		cookie := http.Cookie{Name: "access_token",
			Value:   newAccessToken,
			Expires: time.Now().Add(refreshTokenLife)}
		http.SetCookie(w, &cookie)
		cookie = http.Cookie{Name: "refresh_token",
			Value:   newRefreshToken,
			Expires: time.Now().Add(refreshTokenLife)}
		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusAccepted)

		return
	}
}

func sendMail(email string) error {
	// Sender data.
	from := viper.GetString("email.from")
	password := viper.GetString("email.password")

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := viper.GetString("email.host")
	smtpPort := viper.GetString("email.port")

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil

}
