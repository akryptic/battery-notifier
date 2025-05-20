package cli

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type Options struct {
	Test       bool
	Ntfy       bool
	Reset      bool
	Read       bool
	DryRun     bool
	ConfigPath string
}

func ParseArgs() (*Options, error) {
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
		return nil, err
	}

	return &Options{
		Test:       *testFlag,
		Ntfy:       *ntfyFlag,
		Reset:      *resetFlag,
		Read:       *readFlag,
		DryRun:     *dryRunFlag,
		ConfigPath: configPath,
	}, nil
}
