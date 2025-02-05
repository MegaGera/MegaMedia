package handlers

import (
	"MegaMedia/internal/config"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// squareImage takes an input image and returns a squared image
func squareImage(img image.Image, bgColor color.Color) image.Image {
	// Get the original dimensions
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Determine the square size (max of width and height)
	squareSize := width
	if height > width {
		squareSize = height
	}

	// Create a new square canvas
	square := image.NewRGBA(image.Rect(0, 0, squareSize, squareSize))

	// Fill the canvas with the background color
	draw.Draw(square, square.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Calculate the offset to center the original image
	offsetX := (squareSize - width) / 2
	offsetY := (squareSize - height) / 2

	// Draw the original image onto the square canvas
	draw.Draw(square, image.Rect(offsetX, offsetY, offsetX+width, offsetY+height), img, img.Bounds().Min, draw.Over)

	return square
}

// saveImage saves an image to a file
func saveImage(img image.Image, teamID string, teamDir string, filename string) error {

	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, rename it with "_old_{timestamp}" suffix
		timestamp := time.Now().Format("20060102_150405")
		oldFileName := fmt.Sprintf("team_%s_old_%s.png", teamID, timestamp)
		oldFilePath := path.Join(teamDir, oldFileName)

		// Rename the existing file
		log.Printf("Renaming existing file: %s -> %s\n", filename, oldFilePath)
		if err := os.Rename(filename, oldFilePath); err != nil {
			return fmt.Errorf("failed to rename existing file: %v", err)
		}

		// Send a POST request to megagoal server to notify about the old file
		postURL := config.Cfg.MegagoalServerApiUrl + "/team/set_previous_image"
		payload := fmt.Sprintf(`{"team_id": "%s", "image_title": "%s"}`, teamID, oldFileName)
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

	// Create the output file
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Determine the format based on the file extension
	if strings.HasSuffix(filename, ".png") {
		err = png.Encode(outFile, img)
	} else if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	} else {
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90}) // Default to JPEG
	}

	return err
}

// SquaredImage handles image squared feature
func SquaredImageHandler(w http.ResponseWriter, r *http.Request, teamID string) error {
	// Construct the file path to retrieve the original image
	teamDir := "./static/megagoal/teams/" // Directory path
	newFilePath := path.Join(teamDir, fmt.Sprintf("team_%s.png", teamID))

	// Read the existing file
	file, err := os.Open(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to open existing file: %v", err)
	}
	defer file.Close()

	// Read the file content into a variable
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to read existing file: %v", err)
	}
	log.Printf("After reading file")

	// Process the image to make it square
	bgColor := color.Transparent
	squaredImage := squareImage(img, bgColor)

	log.Printf("After squared")

	// Save the squared image
	// outputFile := path.Join(teamDir, fmt.Sprintf("team_%s_squared.png", teamID))
	err = saveImage(squaredImage, teamID, teamDir, newFilePath)
	if err != nil {
		return fmt.Errorf("error saving image: %v", err)
	}
	log.Printf("After save")

	log.Printf("File uploaded successfully for team %s: %s\n", teamID, newFilePath)
	return nil
}
