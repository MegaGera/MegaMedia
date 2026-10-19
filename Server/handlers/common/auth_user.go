package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"MegaMedia/internal/config"
)

var safeUsernamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,64}$`)

// SanitizeProfileUsername validates a username for use as a profile image filename stem.
func SanitizeProfileUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return "", fmt.Errorf("empty username")
	}
	if !safeUsernamePattern.MatchString(username) {
		return "", fmt.Errorf("invalid username")
	}
	return username, nil
}

type authValidateResponse struct {
	Message string `json:"message"`
	Data    struct {
		Username string `json:"username"`
	} `json:"data"`
}

// UsernameFromAccessToken asks MegaAuth /validate and returns the JWT username.
func UsernameFromAccessToken(token string) (string, error) {
	if strings.TrimSpace(token) == "" {
		return "", fmt.Errorf("missing token")
	}
	if config.Cfg.ValidateUri == "" {
		return "", fmt.Errorf("VALIDATE_URI not configured")
	}

	req, err := http.NewRequest("GET", config.Cfg.ValidateUri, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create validate request: %v", err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: token,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("failed to read validate response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unauthorized")
	}

	var parsed authValidateResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("invalid validate response")
	}
	return SanitizeProfileUsername(parsed.Data.Username)
}
