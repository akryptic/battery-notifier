package main

import (
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/akryptic/battery-notifier/internal/battery"
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/notification"
)

func main() {

	configPath := fmt.Sprintf("%s/.config/battery-notifier/config.toml", os.Getenv("HOME"))

	parser := argparse.NewParser("battery-notifier", "A simple utility that monitors battery levels and sends notifications when the battery is low, critically low, or overcharged.")

	testFlag := parser.Flag("t", "test", &argparse.Options{
		Required: false,
		Help:     "Send test notification",
	})

	ntfyFlag := parser.Flag("n", "ntfy", &argparse.Options{
		Required: false,
		Help:     "Send notification via ntfy",
	})

	resetFlag := parser.Flag("r", "reset", &argparse.Options{
		Required: false,
		Help:     "Reset config to default (use carefully as this will overwrite your current config)",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if *resetFlag {
		config.GenerateDefaultConfig(configPath)
		return
	}

	conf, err := config.Load(configPath)

	validationError := conf.Validate()

	if validationError != "" {
		fmt.Println(validationError)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if *ntfyFlag {

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

	if *testFlag {
		err := notification.SendNotification("Battery Notifier Test", "This is a test notification.", &conf)

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	notifiedLow := false
	notifiedCritical := false
	notifiedOvercharge := false

	for {
		capacity, err := battery.ReadCapacity()
		if err != nil {
			fmt.Println(err)
			continue
		}

		status, err := battery.ReadStatus()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("[%s] Battery: %d%% %s\n", time.Now().Format("2006-01-02 15:04:05"), capacity, status)

		// discharging, under low battery
		if capacity <= conf.LowBattery && !notifiedLow && status != "Charging" {
			err := notification.SendNotification(
				"Battery Low",
				fmt.Sprintf("%d%% remaining. Please plug in.", capacity),
				&conf,
			)
			if err != nil {
				fmt.Println(err)
			}
			notifiedLow = true
		}

		// charging, over low battery
		if capacity > conf.LowBattery && notifiedLow {
			notifiedLow = false
		}

		// discharging, under critical battery
		if capacity <= conf.CriticalBattery && !notifiedCritical && status != "Charging" {
			err := notification.SendNotification(
				"Battery Critically Low",
				fmt.Sprintf("%d%% remaining! System may shut down.", capacity),
				&conf,
			)
			if err != nil {
				fmt.Println(err)
			}
			notifiedCritical = true
		}

		// charging, over critical battery
		if capacity > conf.CriticalBattery && notifiedCritical {
			notifiedCritical = false
		}

		// charging, over overcharge limit
		if capacity >= conf.OverchargeLimit && status == "Charging" && !notifiedOvercharge {
			err := notification.SendNotification(
				"Battery Overcharging",
				fmt.Sprintf(
					"%d%% charged. Consider unplugging to preserve battery health.",
					capacity,
				),
				&conf,
			)
			if err != nil {
				fmt.Println(err)
			}
			notifiedOvercharge = true
		}

		// discharging, under overcharge limit
		if capacity < conf.OverchargeLimit && notifiedOvercharge {
			notifiedOvercharge = false
		}

		time.Sleep(time.Duration(conf.CheckInterval) * time.Second)
	}

}
