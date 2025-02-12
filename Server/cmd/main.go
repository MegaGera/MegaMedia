package main

import (
	"MegaMedia/handlers"
	"log"
	"net/http"
	"strings"

	"MegaMedia/internal/config"
)

func main() {
	// Load environment variables
	config.LoadConfig() // Load config once

	// API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/megagoal/", handlers.MegagoalHandler)
	apiMux.HandleFunc("/api/megagera/", handlers.MegageraHandler)

	// Static file serving (no CORS applied here)
	staticMux := http.NewServeMux()
	staticMux.Handle("/megagoal/", http.StripPrefix("/megagoal/", http.FileServer(http.Dir("./static/megagoal"))))
	staticMux.Handle("/megagera/", http.StripPrefix("/megagera/", http.FileServer(http.Dir("./static/megagera"))))

	// Middleware to validate API requests
	protectedAPI := withValidationMiddleware(apiMux)

	// Apply custom CORS middleware
	corsHandler := customCORSMiddleware(protectedAPI)

	// Create main mux to handle both API & static files
	mainMux := http.NewServeMux()
	mainMux.Handle("/api/megagoal/", corsHandler) // API with CORS
	mainMux.Handle("/api/megagera/", corsHandler)
	mainMux.Handle("/", staticMux) // Static files, no CORS

	// Start server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mainMux))
}

func customCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow all subdomains of megagera.com
		if config.Cfg.AppEnv != "development" && (origin != "" && (strings.HasSuffix(origin, ".megagera.com") || origin == "https://megagera.com")) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if config.Cfg.AppEnv == "development" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight OPTIONS request
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			http.Error(w, "CORS request from disallowed origin", http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}

func withValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only continue if is production
		if config.Cfg.AppEnv != "development" {

			// Extract token from cookie
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "Unauthorized: No token provided main", http.StatusUnauthorized)
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
		} else {
			// Proceed to the next handler
			next.ServeHTTP(w, r)
		}
	})
}

func validateToken(token string) bool {
	req, err := http.NewRequest("GET", config.Cfg.ValidateUri, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return false
	}

	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: token,
	})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request to auth service:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
