package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/CoryEvans2324/eds-enterprise-notes/config"
	"github.com/CoryEvans2324/eds-enterprise-notes/models"
	"github.com/golang-jwt/jwt"
)

const JWT_TOKEN_COOKIE_NAME = "enterprisenotesauth"

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type ContextUserKey struct{}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user and insert it into the request context

		cookie, err := r.Cookie(JWT_TOKEN_COOKIE_NAME)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &models.JWTUser{}, func(t *jwt.Token) (interface{}, error) {
			return config.Get().SecretAsBytes(), nil
		})
		claims, ok := token.Claims.(*models.JWTUser)
		if !ok || !token.Valid || err != nil {
			log.Println("[JWT ERROR]", err.Error())
		}

		ctx := context.WithValue(r.Context(), ContextUserKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(r *http.Request) *models.JWTUser {
	ctx := r.Context()
	user := ctx.Value(ContextUserKey{})
	if user == nil {
		return nil
	}
	userJwt := user.(*models.JWTUser)
	return userJwt
}

func SetUser(w http.ResponseWriter, jwtUser *models.JWTUser) {
	SetJWTCookie(w, jwtUser)
}

func SetJWTCookie(w http.ResponseWriter, jwtUser *models.JWTUser) {
	if jwtUser == nil {
		http.SetCookie(
			w,
			&http.Cookie{
				Name:    JWT_TOKEN_COOKIE_NAME,
				Value:   "",
				Expires: time.Unix(0, 0),
				Path:    "/",
			},
		)
		return
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		jwtUser,
	)

	tokenString, _ := token.SignedString(config.Get().SecretAsBytes())
	http.SetCookie(
		w,
		&http.Cookie{
			Name:    JWT_TOKEN_COOKIE_NAME,
			Value:   tokenString,
			Expires: time.Now().Add(365 * 24 * time.Hour),
			Path:    "/",
		},
	)
}
