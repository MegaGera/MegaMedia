package handlers

import (
	commonHandlers "MegaMedia/handlers/common"
	"fmt"
	"net/http"
	"strings"
)

// megageraHandler handles requests to /api/megagera/{name}/{operationType}
func MegageraHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK) // Preflight requests should return 200 OK
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the fileName from the URL
	operation := strings.TrimPrefix(r.URL.Path, "/api/megagera/")
	operation = strings.TrimSuffix(operation, "/")

	if r.Method == http.MethodGet {
		if operation == "" {
			if err := commonHandlers.FetchImages(w, r, "megagera", "logo"); err != nil {
				http.Error(w, fmt.Sprintf("Error fetching logos: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Invalid get operation", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		if operation == "create" {
			// Handle creating a new image
			if err := commonHandlers.CreateImage(w, r, "megagera", "logo", "megagera"); err != nil {
				http.Error(w, fmt.Sprintf("Error creating image: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}
		operationArray := strings.Split(operation, "/")
		imageID := operationArray[0]
		operationType := operationArray[1]

		switch operationType {
		case "update-name":
			// Handle updating image name
			if err := commonHandlers.UpdateImageName(w, r, "megagera", "logo", imageID); err != nil {
				http.Error(w, fmt.Sprintf("Error updating image name: %v", err), http.StatusInternalServerError)
				return
			}
			return
		case "squared":
			// Handle squaring the image
			image, err := commonHandlers.FetchImageByID(w, r, "megagera", "logo", imageID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching image: %v", err), http.StatusInternalServerError)
				return
			}
			previousFileName, err := commonHandlers.SquaredMegageraImageHandler(w, r, "megagera", image.URL)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error squaring image: %v", err), http.StatusInternalServerError)
				return
			}
			if previousFileName != "" {
				// Update megamedia database (mongodb) adding the previous new file name
				if err := commonHandlers.UpdateImagesPrevious(w, r, "megagera", "logo", image.ID, previousFileName, "add"); err != nil {
					http.Error(w, fmt.Sprintf("Error updating in megamedia mongodb: %v", err), http.StatusInternalServerError)
					return
				}
			}
			return
		case "upload":
			// Handle the image upload
			image, err := commonHandlers.FetchImageByID(w, r, "megagera", "logo", imageID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error uploading image: %v", err), http.StatusInternalServerError)
				return
			}
			nameSplitted := strings.Split(image.URL, ".")
			previousFileName, err := commonHandlers.UploadImageHandler(w, r, "/megagera/", strings.Join(nameSplitted[:len(nameSplitted)-1], "."))
			if err != nil {
				http.Error(w, fmt.Sprintf("Error uploading image: %v", err), http.StatusInternalServerError)
				return
			}
			if previousFileName != "" {
				// Update megamedia database (mongodb) adding the previos new file name
				if err := commonHandlers.UpdateImagesPrevious(w, r, "megagera", "logo", image.ID, previousFileName, "add"); err != nil {
					http.Error(w, fmt.Sprintf("Error updating in megamedia mongodb: %v", err), http.StatusInternalServerError)
					return
				}
			}
		case "delete":
			previousFileName := operationArray[2]
			image, err := commonHandlers.FetchImageByID(w, r, "megagera", "logo", imageID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error uploading image: %v", err), http.StatusInternalServerError)
				return
			}
			if err := commonHandlers.DeleteImageHandler(w, r, "megagera", previousFileName); err != nil {
				http.Error(w, fmt.Sprintf("Error deleting image: %v", err), http.StatusInternalServerError)
				return
			}
			if previousFileName != "" {
				// Update megamedia database (mongodb) deleting the previos new file name
				if err := commonHandlers.UpdateImagesPrevious(w, r, "megagera", "logo", image.ID, previousFileName, "delete"); err != nil {
					http.Error(w, fmt.Sprintf("Error deleting in megamedia mongodb: %v", err), http.StatusInternalServerError)
					return
				}
			}
		default:
			http.NotFound(w, r)
			return
		}
	}

}
