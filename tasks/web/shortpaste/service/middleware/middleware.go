package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"log"
	"net/http"

	"shortpaste/config"
	"shortpaste/models"
)

type UserAuth struct {
	User         models.User
	IsAuthorized bool
}

type AuthMiddleware struct {
	DB *gorm.DB
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := UserAuth{IsAuthorized: false}

		authCookie, err := r.Cookie("auth")
		if err != nil {
			ctx := context.WithValue(r.Context(), "userAuth", ua)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwt.Parse(authCookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.HMAC_SECRET), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if user_id, ok := claims["id"]; ok {
				var user models.User
				tx := am.DB.Model(&user).Preload("Pastes").Find(&user, "id = ?", user_id)
				if tx.RowsAffected != 0 {
					ua.IsAuthorized = true
					ua.User = user
				}
			}
		} else {
			log.Print(err)
		}

		ctx := context.WithValue(r.Context(), "userAuth", ua)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
