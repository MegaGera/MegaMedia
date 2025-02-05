package main

import (
	"MegaMedia/handlers"
	"fmt"
	"log"
	"net/http"
	"strings"

	"MegaMedia/internal/config"

	"github.com/rs/cors"
)

func main() {
	// Load environment variables
	config.LoadConfig() // Load config once

	// API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/megagoal/team/", teamHandler)

	// Static file serving (no CORS applied here)
	staticMux := http.NewServeMux()
	staticMux.Handle("/megagoal/", http.StripPrefix("/megagoal/", http.FileServer(http.Dir("./static/megagoal"))))
	staticMux.Handle("/megagera/", http.StripPrefix("/megagera/", http.FileServer(http.Dir("./static/megagera"))))

	// Middleware to validate API requests
	protectedAPI := withValidationMiddleware(apiMux)

	// Apply CORS only to API
	var allowedOrigins []string
	if config.Cfg.AppEnv == "development" {
		allowedOrigins = []string{"*"}
	} else {
		allowedOrigins = []string{`/\.?megagera\.com$/`}
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(protectedAPI)

	// Create main mux to handle both API & static files
	mainMux := http.NewServeMux()
	mainMux.Handle("/megagoal/team/", corsHandler) // API with CORS
	mainMux.Handle("/", staticMux)                 // Static files, no CORS

	// Start server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mainMux))
}

func withValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}
		token := cookie.Value

		// Validate token with authentication microservice
		if !validateToken(token) {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if valid
		next.ServeHTTP(w, r)
	})
}

func validateToken(token string) bool {
	req, err := http.NewRequest("GET", config.Cfg.AppEnv, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return false
	}

	req.Header.Set("access_token", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request to auth service:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// teamHandler handles requests to /team/{team_id}/image
func teamHandler(w http.ResponseWriter, r *http.Request) {
	// Normalize path by removing trailing slash
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the teamID from the URL
	operation := strings.TrimPrefix(r.URL.Path, "/megagoal/team/")
	operationArray := strings.Split(operation, "/")

	if len(operationArray) < 2 {
		http.NotFound(w, r)
		return
	}

	teamID := operationArray[0]
	operationType := operationArray[1]

	switch operationType {
	case "squared":
		// Handle the image squared
		log.Printf("Squared image")
		if err := handlers.SquaredImageHandler(w, r, teamID); err != nil {
			http.Error(w, fmt.Sprintf("Error squaring image: %v", err), http.StatusInternalServerError)
			return
		}
	case "image":
		// Handle the image upload
		log.Printf("Upload image")
		if err := handlers.UploadImageHandler(w, r, teamID); err != nil {
			http.Error(w, fmt.Sprintf("Error uploading image: %v", err), http.StatusInternalServerError)
			return
		}
	case "delete":
		// Handle the image deleted
		log.Printf("Delete image")
		if err := handlers.DeleteImageHandler(w, r, teamID); err != nil {
			http.Error(w, fmt.Sprintf("Error deleting image: %v", err), http.StatusInternalServerError)
			return
		}
	default:
		http.NotFound(w, r)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Image uploaded successfully for team %s", teamID)
}
