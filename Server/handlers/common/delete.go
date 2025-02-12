package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

// DeleteImageHandler handles image deletion for a specific team and returns an error if something goes wrong
func DeleteImageHandler(w http.ResponseWriter, r *http.Request, staticRoute string, fileName string) error {

	// Create the filepath to delete the image
	fileDir := path.Join("./static", staticRoute, fileName)

	// Check if the file already exists
	if _, err := os.Stat(fileDir); err == nil {

		// Delete the existing file
		log.Printf("Deleting existing file: %s\n", fileDir)
		if err := os.Remove(fileDir); err != nil {
			return fmt.Errorf("failed to delete existing file: %v", err)
		}

	} else {
		return fmt.Errorf("file does not exist")
	}

	log.Printf("File deleted successfully for team: %s\n", fileName)
	return nil
}
