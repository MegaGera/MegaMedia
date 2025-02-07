package handlers

import (
	"MegaMedia/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// DeleteImageHandler handles image deletion for a specific team and returns an error if something goes wrong
func DeleteImageHandler(w http.ResponseWriter, r *http.Request, teamID string) error {

	// Create the filepath to delete the image
	teamDir := "./static/megagoal/teams/" // Directory path
	newFilePath := path.Join(teamDir, teamID)

	// Split the teamID by "_" and take the second position
	parts := strings.Split(teamID, "_")
	if len(parts) < 2 {
		return fmt.Errorf("invalid teamID format")
	}
	teamIDReal := parts[1]

	// Check if the file already exists
	if _, err := os.Stat(newFilePath); err == nil {

		// Delete the existing file
		log.Printf("Deleting existing file: %s\n", newFilePath)
		if err := os.Remove(newFilePath); err != nil {
			return fmt.Errorf("failed to delete existing file: %v", err)
		}

		// Send a POST request to megagoal server to notify about the deletion the old file
		postURL := config.Cfg.MegagoalServerApiUrl + "/team/delete_previous_image"
		payload := fmt.Sprintf(`{"team_id": "%s", "image_title": "%s"}`, teamIDReal, teamID)
		req, err := http.NewRequest("POST", postURL, strings.NewReader(payload))

		if err != nil {
			return fmt.Errorf("failed to create POST request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		if config.Cfg.AppEnv != "development" {
			// Extract token from cookie
			cookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "Unauthorized: No token provided in delete", http.StatusUnauthorized)
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
	} else {
		return fmt.Errorf("file does not exist")
	}

	log.Printf("File deleted successfully for team %s: %s\n", teamIDReal, teamID)
	return nil
}
