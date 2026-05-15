package handlers

import (
	commonHandlers "MegaMedia/handlers/common"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// OttoHandler handles POST /api/otto/upload/ with multipart fields:
//   - image (file)
//   - name (filename stem only; extension comes from the uploaded file)
func OttoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	if path != "/api/otto/upload" {
		http.NotFound(w, r)
		return
	}

	urlPath, err := commonHandlers.UploadOttoImage(r)
	if err != nil {
		status := http.StatusInternalServerError
		msg := err.Error()
		switch msg {
		case "empty name", "name too long", "invalid name", "file has no extension":
			status = http.StatusBadRequest
		default:
			if strings.HasPrefix(msg, "failed to retrieve file:") {
				status = http.StatusBadRequest
			}
		}
		http.Error(w, fmt.Sprintf("Error uploading image: %v", err), status)
		return
	}

	body, err := json.Marshal(map[string]string{"url": urlPath})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
