package notification

import (
	"fmt"
	"os"
	"os/exec"

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

	if config.EnableSound {
		soundPath := config.SoundFile
		if soundPath == "" {
			soundPath = "bell"
		}

		volume := config.SoundVolume
		if volume < 0 {
			volume = 0
		}
		if volume > 100 {
			volume = 100
		}

		// TODO: cross platform support
		cmd := exec.Command("canberra-gtk-play", "-f", soundPath, fmt.Sprintf("--volume=%d", volume))

		err := cmd.Start()
		if err != nil {
			return err
		}

		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}
