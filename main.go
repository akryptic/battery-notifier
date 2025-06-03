package main

import (
	"fmt"

	"github.com/akryptic/battery-notifier/internal/battery"
	"github.com/akryptic/battery-notifier/internal/cli"
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/logging"
	"github.com/akryptic/battery-notifier/internal/monitor"
	"github.com/akryptic/battery-notifier/internal/notification"
	"github.com/akryptic/battery-notifier/internal/sound"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Print(err)
		return
	}

	logging.Info("Starting battery-notifier with verbose level %d", opts.VerboseLevel)

	logging.Debug("Parsed args: %+v", opts)

	// --reset
	if opts.Reset {
		logging.Trace("Detected --reset flag")
		config.GenerateDefaultConfig(opts.ConfigPath)
		return
	}

	// --read
	if opts.Read {
		logging.Trace("Detected --read flag")

		state, err := battery.ReadBatteryState()

		if err != nil {
			logging.Error("Failed to read battery state: %v", err)
			return
		}

		logging.Info("Status: %s", state.Status)
		logging.Info("Level: %d%%", state.Level)
		return
	}

	conf, err := config.Load(opts.ConfigPath)
	if err != nil {
		logging.Error("Failed to load config: %v", err)
		return
	}

	validationError := conf.Validate()
	if validationError != "" {
		logging.Error("Failed to validate config: %s", validationError)
		return
	}

	// --ntfy
	if opts.Ntfy {
		logging.Trace("Detected --ntfy flag")

		if conf.NtfyTopic == "" {
			logging.Error("No ntfy topic set in config.")
			return
		}

		err := notification.SendNtfyNotification("This is a test notification.", conf.GetNtfyUrl(), conf.NtfyAccessToken)
		if err != nil {
			logging.Error("Failed to send ntfy notification: %v", err)
			return
		}

		logging.Info("Test ntfy notification sent successfully")
		return
	}

	// --test
	if opts.Test {
		logging.Trace("Detected --test flag")

		conf.NtfyTopic = ""
		err := notification.SendNotification("Battery Notifier Test", "This is a test notification.", &conf)
		if err != nil {
			logging.Error("Failed to send test notification: %v", err)
			return
		}

		if conf.EnableSound {
			err := sound.Play("low", &conf)
			if err != nil {
				logging.Error("Failed to play sound for test notification: %v", err)
			}
		}
		logging.Info("Test notification sent successfully")
		return
	}

	batteryMonitor := monitor.NewMonitor(&conf)

	// --dry-run
	if opts.DryRun {
		logging.Trace("Detected --dry-run flag")

		runType, err := batteryMonitor.RunOnce()
		if err != nil {
			logging.Error("Dry run failed: %v", err)
			return
		}

		logging.Info("Dry run: %s", runType)
		return
	}

	batteryMonitor.StartMonitoring()
}
