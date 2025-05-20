package notification

import (
	"net/http"
	"strings"
)

func SendNtfyNotification(message string, ntfyEndpoint string, ntfyAccessToken string) error {
	if ntfyEndpoint == "" {
		return nil
	}

	req, err := http.NewRequest("POST", ntfyEndpoint, strings.NewReader(message))

	if err != nil {
		return err
	}

	// attach access token if provided
	if ntfyAccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+ntfyAccessToken)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
