package notification

import (
	"fmt"
	"os"

	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/gen2brain/beeep"
)

func SendNotification(title, message string, config *config.Config) error {
	err := beeep.Notify(title, message, "")
	if err != nil {
		return err
	}

	if config.NtfyTopic != "" {
		ntfyErr := SendNtfyNotification(message, config.GetNtfyUrl(), config.NtfyAccessToken)

		if ntfyErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to send ntfy notification: %v\n", ntfyErr)
		}
	}

	return nil
}
