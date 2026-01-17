package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SigninHandler(cfg *APIConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, fmt.Errorf("invalid json format"), http.StatusBadRequest)
			return
		}

		if cfg.TodoPassword == "" || body.Password != cfg.TodoPassword {
			writeError(w, fmt.Errorf("invalid password"), http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"pwd": body.Password,
			"exp": time.Now().Add(cfg.TokenExpiry).Unix(),
		})

		tokenString, err := token.SignedString([]byte(cfg.TodoPassword))
		if err != nil {
			writeError(w, err, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   int(cfg.TokenExpiry.Seconds()),
		})

		writeJSON(w, map[string]string{"token": tokenString})
	}
}

func Auth(cfg *APIConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.TodoPassword == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeError(w, fmt.Errorf("authentication required"), http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(cfg.TodoPassword), nil
		})

		if err != nil || !token.Valid {
			writeError(w, fmt.Errorf("authentication required"), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if pwd, ok := claims["pwd"].(string); !ok || pwd != cfg.TodoPassword {
				writeError(w, fmt.Errorf("authentication required"), http.StatusUnauthorized)
				return
			}
		} else {
			writeError(w, fmt.Errorf("authentication required"), http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
