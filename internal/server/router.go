package server

import (
	"net/http"
	"path/filepath"
	"time"

	"Final_Project/internal/api"
	"Final_Project/internal/settings"
)

func NewRouter(cfg settings.Config) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(cfg.WebDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(cfg.WebDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})

	apiCfg := &api.APIConfig{
		TodoPassword: cfg.TODO_PASSWORD,
		TokenExpiry:  8 * time.Hour,
	}

	api.Init(mux, apiCfg)

	return mux
}
