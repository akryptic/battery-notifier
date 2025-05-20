package main

import (
	"fmt"

	"github.com/akryptic/battery-notifier/internal/battery"
	"github.com/akryptic/battery-notifier/internal/cli"
	"github.com/akryptic/battery-notifier/internal/config"
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

	// --reset
	if opts.Reset {
		config.GenerateDefaultConfig(opts.ConfigPath)
		return
	}

	// --read
	if opts.Read {
		state, err := battery.ReadBatteryState()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Status: %s\n", state.Status)
		fmt.Printf("Level: %d%%\n", state.Level)
		return
	}

	conf, err := config.Load(opts.ConfigPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	validationError := conf.Validate()
	if validationError != "" {
		fmt.Println(validationError)
		return
	}

	// --ntfy
	if opts.Ntfy {
		if conf.NtfyTopic == "" {
			fmt.Println("No ntfy topic set in config.")
			return
		}

		err := notification.SendNtfyNotification("This is a test notification.", conf.GetNtfyUrl(), conf.NtfyAccessToken)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// --test
	if opts.Test {
		conf.NtfyTopic = ""
		err := notification.SendNotification("Battery Notifier Test", "This is a test notification.", &conf)
		if err != nil {
			fmt.Println(err)
		}

		if conf.EnableSound {
			err := sound.Play("low", &conf)
			if err != nil {
				fmt.Println(err)
			}
		}
		return
	}

	batteryMonitor := monitor.NewMonitor(&conf)

	// --dry-run
	if opts.DryRun {
		runType, err := batteryMonitor.RunOnce()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Dry run: %s\n", runType)
		return
	}

	batteryMonitor.StartMonitoring()
}
