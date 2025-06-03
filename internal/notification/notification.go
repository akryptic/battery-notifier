package notification

import (
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/logging"
	"github.com/gen2brain/beeep"
)

func SendNotification(title, message string, config *config.Config) error {
	logging.Debug("Sending notification: title='%s', message='%s'", title, message)

	err := beeep.Notify(title, message, "")
	if err != nil {
		logging.Error("Failed to send local notification: %v", err)
		return err
	}

	logging.Debug("Local notification sent successfully")

	if config.NtfyTopic != "" {
		ntfyErr := SendNtfyNotification(message, config.GetNtfyUrl(), config.NtfyAccessToken)

		if ntfyErr != nil {
			logging.Error("Failed to send ntfy notification: %v", ntfyErr)
		}
	}

	return nil
}
