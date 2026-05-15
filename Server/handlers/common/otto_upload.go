package handlers

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

// Otto uploads: 10 MiB multipart body cap.
const ottoMaxMultipartBytes = int64(10 << 20)

const ottoUploadAttempts = 16

// SanitizeOttoBaseName validates the logical filename stem from the API (no path segments).
func SanitizeOttoBaseName(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("empty name")
	}
	if len(s) > 200 {
		return "", fmt.Errorf("name too long")
	}
	if strings.Contains(s, "..") || strings.Contains(s, "/") || strings.Contains(s, "\\") {
		return "", fmt.Errorf("invalid name")
	}
	return s, nil
}

func randomSixDigitSuffix() (string, error) {
	var buf [4]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	n := binary.BigEndian.Uint32(buf[:]) % 1_000_000
	return fmt.Sprintf("%06d", n), nil
}

// UploadOttoImage reads multipart fields "image" (file) and "name" (filename stem).
// Files are stored under ./static/otto/ as name_XXXXXX.ext (ext from the uploaded file).
// Returns the public URL path (e.g. /otto/foo_012345.jpg).
func UploadOttoImage(r *http.Request) (publicURLPath string, err error) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("UploadOttoImage panic: %v\n%s", p, debug.Stack())
			publicURLPath = ""
			err = fmt.Errorf("internal upload error")
		}
	}()

	if err := r.ParseMultipartForm(ottoMaxMultipartBytes); err != nil {
		return "", fmt.Errorf("failed to parse form: %v", err)
	}

	rawName := strings.TrimSpace(r.FormValue("name"))
	rawName = path.Base(rawName)
	if ext := path.Ext(rawName); ext != "" {
		rawName = strings.TrimSuffix(rawName, ext)
	}

	baseName, err := SanitizeOttoBaseName(rawName)
	if err != nil {
		return "", err
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file: %v", err)
	}
	if file == nil || handler == nil {
		return "", fmt.Errorf("failed to retrieve file: missing file or header")
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	ext := path.Ext(handler.Filename)
	if ext == "" {
		return "", fmt.Errorf("file has no extension")
	}

	// Same layout as megagoal/megagera: ./static/<product>/...
	fileDir := path.Join(".", "static", "otto")
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	for range ottoUploadAttempts {
		suffix, err := randomSixDigitSuffix()
		if err != nil {
			return "", fmt.Errorf("failed to generate suffix: %v", err)
		}
		stem := fmt.Sprintf("%s_%s", baseName, suffix)
		newFilePath := path.Join(fileDir, stem+ext)

		if _, err := os.Stat(newFilePath); err == nil {
			continue // rare collision; retry (file not read yet)
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to stat file: %v", err)
		}

		log.Printf("Saving Otto image to: %s\n", newFilePath)
		destFile, err := os.Create(newFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to create file: %v", err)
		}
		_, copyErr := io.Copy(destFile, file)
		_ = destFile.Close()
		if copyErr != nil {
			_ = os.Remove(newFilePath)
			return "", fmt.Errorf("failed to save file: %v", copyErr)
		}

		return path.Join("/", "otto", stem+ext), nil
	}

	return "", fmt.Errorf("could not allocate a unique filename after %d attempts", ottoUploadAttempts)
}
