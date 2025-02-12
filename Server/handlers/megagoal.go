package handlers

import (
	commonHandlers "MegaMedia/handlers/common"
	"MegaMedia/internal/config"
	"fmt"
	"net/http"
	"strings"
)

func megagoalServerUpdate(w http.ResponseWriter, r *http.Request, teamID string, newFileName string, operation string) error {
	// Send a POST request to megagoal server to notify about the old file
	var postURL string
	if operation == "add" {
		postURL = config.Cfg.MegagoalServerApiUrl + "/team/set_previous_image"
	} else if operation == "delete" {
		postURL = config.Cfg.MegagoalServerApiUrl + "/team/delete_previous_image"
	}
	payload := fmt.Sprintf(`{"team_id": "%s", "image_title": "%s"}`, teamID, newFileName)
	req, err := http.NewRequest("POST", postURL, strings.NewReader(payload))

	if err != nil {
		return fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add JWT token for MegaGoal Server operation
	if config.Cfg.AppEnv != "development" {
		// Extract token from cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized: No token provided for megagoal", http.StatusUnauthorized)
			return nil
		}
		token := cookie.Value

		req.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: token,
		})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from server: %s", resp.Status)
	}
	return nil
}

// MegagoalHandler handles requests to /api/megagoal/{collection}/{teamID}/{operationType}
func MegagoalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK) // Preflight requests should return 200 OK
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the teamID from the URL
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	operation := strings.TrimPrefix(r.URL.Path, "/api/megagoal/")
	operationArray := strings.Split(operation, "/")

	if len(operationArray) < 3 {
		http.NotFound(w, r)
		return
	}

	collection := operationArray[0]
	teamID := operationArray[1]
	operationType := operationArray[2]

	if collection == "team" {
		switch operationType {
		case "squared":
			// Handle the image squared
			previousFileName, err := commonHandlers.SquaredImageHandler(w, r, "/megagoal/teams/", "team_"+teamID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error squaring image: %v", err), http.StatusInternalServerError)
				return
			}
			if previousFileName != "" {
				// Send a POST request to megagoal server to notify about update old file
				if err := megagoalServerUpdate(w, r, teamID, previousFileName, "add"); err != nil {
					http.Error(w, fmt.Sprintf("Error sending POST request to megagoal: %v", err), http.StatusInternalServerError)
					return
				}
			}
		case "image":
			// Handle the image upload
			previousFileName, err := commonHandlers.UploadImageHandler(w, r, "/megagoal/teams/", "team_"+teamID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error uploading image: %v", err), http.StatusInternalServerError)
				return
			}
			if previousFileName != "" {
				// Send a POST request to megagoal server to notify about update old file
				if err := megagoalServerUpdate(w, r, teamID, previousFileName, "add"); err != nil {
					http.Error(w, fmt.Sprintf("Error sending POST request to megagoal: %v", err), http.StatusInternalServerError)
					return
				}
			}
		case "delete":
			// Handle the image deleted
			previousFileName := teamID
			// Split the teamID by "_" and take the second position
			parts := strings.Split(teamID, "_")
			if len(parts) < 2 {
				http.Error(w, "invalid previousFileName / teamID format", http.StatusInternalServerError)
			}
			teamID := parts[1]

			if err := commonHandlers.DeleteImageHandler(w, r, "/megagoal/teams/", previousFileName); err != nil {
				http.Error(w, fmt.Sprintf("Error deleting image: %v", err), http.StatusInternalServerError)
				return
			}
			// Send a POST request to megagoal server to notify about deletion of old file
			if err := megagoalServerUpdate(w, r, teamID, previousFileName, "delete"); err != nil {
				http.Error(w, fmt.Sprintf("Error sending POST request to megagoal: %v", err), http.StatusInternalServerError)
				return
			}
		default:
			http.NotFound(w, r)
			return
		}

		// Success response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Image uploaded successfully for team %s", teamID)
	} else {
		http.NotFound(w, r)
	}
}
