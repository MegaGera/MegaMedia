package handlers

import (
	"net/http"
	"path/filepath"
)

func ServeImageHandler(w http.ResponseWriter, r *http.Request) {
	imagePath := filepath.Join("./images", r.URL.Path)
	http.ServeFile(w, r, imagePath)
}