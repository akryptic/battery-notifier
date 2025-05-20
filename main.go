package main

import (
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
	"github.com/akryptic/battery-notifier/internal/battery"
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/notification"
	"github.com/akryptic/battery-notifier/internal/sound"
)

type RunType = string

const (
	NoRun              RunType = "No Run"
	LowBatteryRun      RunType = "Low Battery Run"
	CriticalBatteryRun RunType = "Critical Battery Run"
	OverchargeLimitRun RunType = "Overcharging Limit Run"
)

type NotifiedState struct {
	Low        bool
	Critical   bool
	Overcharge bool
}

func determineRunType(conf *config.Config) (RunType, battery.BatteryState, error) {

	batteryState, err := battery.ReadBatteryState()

	if err != nil {
		return NoRun, batteryState, err
	}

	fmt.Printf("[%s] Battery: %d%% %s\n", time.Now().Format("2006-01-02 15:04:05"), batteryState.Level, batteryState.Status)

	if batteryState.Level <= conf.LowBattery && batteryState.Status != battery.Charging {
		return LowBatteryRun, batteryState, nil
	}

	if batteryState.Level <= conf.CriticalBattery && batteryState.Status != battery.Charging {
		return CriticalBatteryRun, batteryState, nil
	}

	if batteryState.Level >= conf.OverchargeLimit && batteryState.Status == battery.Charging {
		return OverchargeLimitRun, batteryState, nil
	}

	return NoRun, batteryState, nil
}

func dryRun(conf *config.Config, notifiedState *NotifiedState) (RunType, error) {

	runType, batteryState, err := determineRunType(conf)

	if err != nil {
		return NoRun, err
	}

	if runType == NoRun {
		return NoRun, nil
	}

	fmt.Println(notifiedState)

	if runType == LowBatteryRun && !notifiedState.Low {
		err := notification.SendNotification(
			"Battery Low",
			fmt.Sprintf("%d%% remaining. Please plug in.", batteryState.Level),
			conf,
		)
		if err != nil {
			return NoRun, err
		}
		if conf.EnableSound {
			err := sound.Play("low", conf)
			if err != nil {
				return NoRun, err
			}
		}

		notifiedState.Low = true

		return LowBatteryRun, nil
	}

	if batteryState.Level > conf.LowBattery && notifiedState.Low {
		notifiedState.Low = false
		return NoRun, nil
	}

	if runType == CriticalBatteryRun && !notifiedState.Critical {
		err := notification.SendNotification(
			"Battery Critically Low",
			fmt.Sprintf("%d%% remaining! System may shut down.", batteryState.Level),
			conf,
		)
		if err != nil {
			return NoRun, err
		}
		if conf.EnableSound {
			err := sound.Play("low", conf)
			if err != nil {
				return NoRun, nil
			}
		}

		notifiedState.Critical = true

		return CriticalBatteryRun, nil
	}

	if batteryState.Level > conf.CriticalBattery && notifiedState.Critical {
		notifiedState.Critical = false
		return NoRun, nil
	}

	if runType == OverchargeLimitRun && !notifiedState.Overcharge {
		err := notification.SendNotification(
			"Battery Overcharging",
			fmt.Sprintf(
				"%d%% charged. Consider unplugging to preserve battery health.",
				batteryState.Level,
			),
			conf,
		)
		if err != nil {
			return NoRun, nil
		}
		if conf.EnableSound {
			err := sound.Play("overcharge", conf)
			if err != nil {
				return NoRun, nil
			}
		}

		notifiedState.Overcharge = true

		return OverchargeLimitRun, nil
	}

	if batteryState.Level < conf.OverchargeLimit && notifiedState.Overcharge {
		notifiedState.Overcharge = false
		return NoRun, nil
	}

	return NoRun, nil
}

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

	readFlag := parser.Flag("", "read", &argparse.Options{
		Required: false,
		Help:     "Reads the battery status and prints it",
	})

	dryRunFlag := parser.Flag("d", "dry-run", &argparse.Options{
		Required: false,
		Help:     "Dry run the notifier",
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

	if *readFlag {
		status, err := battery.ReadStatus()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Status: %s\n", status)

		level, err := battery.ReadLevel()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Level: %d%%\n", level)
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

	notifiedState := NotifiedState{false, false, false}

	if *dryRunFlag {
		runType, err := dryRun(&conf, &notifiedState)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Dry run: %s\n", runType)
		return
	}

	for {
		runType, err := dryRun(&conf, &notifiedState)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Run: %s\n", runType)
		time.Sleep(time.Duration(conf.CheckInterval) * time.Second)
	}

}
