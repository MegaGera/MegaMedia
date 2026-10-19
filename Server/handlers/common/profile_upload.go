package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const profileMaxMultipartBytes = int64(10 << 20)

var allowedProfileExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

// UploadUserProfileImage stores a profile picture under static/megagoal/users/{username}{ext}.
// Returns relative path under megagoal (e.g. "users/gera.png").
func UploadUserProfileImage(r *http.Request, username string) (relativePath string, err error) {
	username, err = SanitizeProfileUsername(username)
	if err != nil {
		return "", err
	}

	// Form may already be parsed by the caller (e.g. development username resolution).
	if r.MultipartForm == nil {
		if err = r.ParseMultipartForm(profileMaxMultipartBytes); err != nil {
			return "", fmt.Errorf("failed to parse form: %v", err)
		}
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file: %v", err)
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(handler.Filename))
	if ext == "" {
		return "", fmt.Errorf("file has no extension")
	}
	if !allowedProfileExtensions[ext] {
		return "", fmt.Errorf("unsupported file type")
	}

	fileDir := path.Join("./static", "megagoal", "users")
	if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Archive any previous profile files for this username (any extension).
	matches, _ := filepath.Glob(path.Join(fileDir, username+".*"))
	timestamp := time.Now().Format("20060102_150405")
	for _, existing := range matches {
		base := path.Base(existing)
		if strings.Contains(base, "_old_") {
			continue
		}
		oldName := fmt.Sprintf("%s_old_%s%s", username, timestamp, path.Ext(existing))
		oldPath := path.Join(fileDir, oldName)
		log.Printf("Renaming existing profile file: %s -> %s\n", existing, oldPath)
		if renameErr := os.Rename(existing, oldPath); renameErr != nil {
			return "", fmt.Errorf("failed to rename existing file: %v", renameErr)
		}
	}

	newFilePath := path.Join(fileDir, username+ext)
	log.Printf("Saving profile file to: %s\n", newFilePath)
	destFile, err := os.Create(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return "users/" + username + ext, nil
}
