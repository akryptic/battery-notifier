package notification

import (
	"net/http"
	"strings"

	"github.com/akryptic/battery-notifier/internal/logging"
)

func SendNtfyNotification(message string, ntfyEndpoint string, ntfyAccessToken string) error {
	logging.Debug("Sending ntfy notification")
	logging.Trace("Ntfy message: '%s'", message)

	if ntfyEndpoint == "" {
		logging.Trace("Empty ntfy endpoint, skipping notification")
		return nil
	}

	req, err := http.NewRequest("POST", ntfyEndpoint, strings.NewReader(message))

	if err != nil {
		logging.Error("Failed to create ntfy HTTP request: %v", err)
		return err
	}

	// attach access token if provided
	if ntfyAccessToken != "" {
		logging.Trace("Adding authorization header to ntfy request")
		req.Header.Set("Authorization", "Bearer "+ntfyAccessToken)
	}

	logging.Trace("Sending ntfy HTTP request")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		logging.Error("Failed to send ntfy HTTP request: %v", err)
		return err
	}

	defer resp.Body.Close()

	logging.Debug("Ntfy HTTP request completed with status: %d", resp.StatusCode)

	return nil
}
