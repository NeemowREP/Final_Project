package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, fmt.Errorf("method not allowed"))
		return
	}

	var body struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, fmt.Errorf("invalid json format"))
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" || body.Password != pass {
		writeError(w, fmt.Errorf("invalid password"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pwd": body.Password,
		"exp": time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(pass))
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   8 * 3600,
	})

	writeJSON(w, map[string]string{"token": tokenString})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeError(w, fmt.Errorf("authentication required"))
			return
		}

		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(pass), nil
		})

		if err != nil || !token.Valid {
			writeError(w, fmt.Errorf("authentication required"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if pwd, ok := claims["pwd"].(string); !ok || pwd != pass {
				writeError(w, fmt.Errorf("authentication required"))
				return
			}
		} else {
			writeError(w, fmt.Errorf("authentication required"))
			return
		}

		next(w, r)
	}
}
