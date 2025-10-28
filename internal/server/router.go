package server

import (
	"Final_Project/internal/api"
	"Final_Project/internal/settings"
	"net/http"
	"path/filepath"
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

	api.Init(mux)

	return mux
}
