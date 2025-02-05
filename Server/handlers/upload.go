package handlers

import (
	"MegaMedia/internal/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// UploadImageHandler handles image uploads for a specific team and returns an error if something goes wrong
func UploadImageHandler(w http.ResponseWriter, r *http.Request, teamID string) error {
	// Parse the incoming multipart form (image upload)
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Get the file from the form (with the key "image")
	file, handler, err := r.FormFile("image")
	if err != nil {
		return fmt.Errorf("failed to retrieve file: %v", err)
	}
	defer file.Close()

	// Get the file extension (e.g., ".jpg", ".png")
	fileExtension := path.Ext(handler.Filename)
	if fileExtension == "" {
		return fmt.Errorf("file has no extension")
	}

	// Create the directory where we will store the uploaded image
	teamDir := "./static/megagoal/teams/" // Directory path
	if err := os.MkdirAll(teamDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create team directory: %v", err)
	}

	// Construct the file path to store the new image
	newFilePath := path.Join(teamDir, fmt.Sprintf("team_%s%s", teamID, fileExtension))

	// Check if the file already exists
	if _, err := os.Stat(newFilePath); err == nil {
		// File exists, rename it with "_old_{timestamp}" suffix
		timestamp := time.Now().Format("20060102_150405")
		newFileName := fmt.Sprintf("team_%s_old_%s%s", teamID, timestamp, fileExtension)
		oldFilePath := path.Join(teamDir, newFileName)

		// Rename the existing file
		log.Printf("Renaming existing file: %s -> %s\n", newFilePath, oldFilePath)
		if err := os.Rename(newFilePath, oldFilePath); err != nil {
			return fmt.Errorf("failed to rename existing file: %v", err)
		}

		// Send a POST request to megagoal server to notify about the old file
		postURL := config.Cfg.MegagoalServerApiUrl + "/team/set_previous_image"
		payload := fmt.Sprintf(`{"team_id": "%s", "image_title": "%s"}`, teamID, newFileName)
		req, err := http.NewRequest("POST", postURL, strings.NewReader(payload))

		if err != nil {
			return fmt.Errorf("failed to create POST request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send POST request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("received non-OK response from server: %s", resp.Status)
		}
	}

	// Create the file on disk for the new image
	log.Printf("Saving new file to: %s\n", newFilePath)
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	// Copy the uploaded file into the newly created file
	_, err = io.Copy(destFile, file)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	log.Printf("File uploaded successfully for team %s: %s\n", teamID, newFilePath)
	return nil
}
