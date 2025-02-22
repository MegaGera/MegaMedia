package handlers

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
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

// SquaredImage handles image squared feature
func SquaredImageHandler(w http.ResponseWriter, r *http.Request, staticRoute string, fileName string) (previousFileName string, err error) {
	// Construct the file path to retrieve the original image
	fileDir := path.Join("./static", staticRoute)
	newFilePath := path.Join(fileDir, fmt.Sprintf("%s.png", fileName))

	// Read the existing file
	file, err := os.Open(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open existing file: %v", err)
	}
	defer file.Close()

	// Read the file content into a variable
	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to read existing file: %v", err)
	}

	// Process the image to make it square
	bgColor := color.Transparent
	squaredImage := squareImage(img, bgColor)

	// File exists, rename it with "_old_{timestamp}" suffix
	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("%s_old_%s.png", fileName, timestamp)
	oldFilePath := path.Join(fileDir, newFileName)

	// Rename the existing file
	log.Printf("Renaming existing file: %s -> %s\n", newFilePath, oldFilePath)
	if err := os.Rename(newFilePath, oldFilePath); err != nil {
		return "", fmt.Errorf("failed to rename existing file: %v", err)
	}

	// Create the file on disk for the new image
	log.Printf("Saving new file to: %s\n", newFilePath)
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	if err = png.Encode(destFile, squaredImage); err != nil {
		return "", fmt.Errorf("failed to save / encode file: %v", err)
	}

	log.Printf("File uploaded successfully for team %s: %s\n", fileName, newFilePath)
	return newFileName, nil
}
