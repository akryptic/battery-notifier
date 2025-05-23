package cli

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/akryptic/battery-notifier/internal/config"
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
	parser := argparse.NewParser("battery-notifier", "A simple utility that monitors battery levels and sends notifications when the battery is low, critically low, or overcharged.")

	testFlag := parser.Flag("t", "test", &argparse.Options{
		Required: false,
		Help:     "Send test notification",
	})

	ntfyFlag := parser.Flag("n", "ntfy", &argparse.Options{
		Required: false,
		Help:     "Send notification via ntfy",
	})

	resetFlag := parser.Flag("", "reset", &argparse.Options{
		Required: false,
		Help:     "Reset config to default (use carefully as this will overwrite your current config)",
	})

	readFlag := parser.Flag("R", "read", &argparse.Options{
		Required: false,
		Help:     "Reads the battery status and prints it",
	})

	dryRunFlag := parser.Flag("d", "dry-run", &argparse.Options{
		Required: false,
		Help:     "Dry run the notifier",
	})

	configPath := parser.String("c", "config", &argparse.Options{
		Required: false,
		Help:     "Path to the config file (default: <config_dir>/battery-notifier/config.toml)",
		Validate: func(args []string) error {
			_, err := config.Load(args[0])
			return err
		},
	})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, err
	}

	// if config path was not passed, use the default path
	if *configPath == "" {
		defaultConfigPath, err := config.GetDefaultConfigPath()

		configPath = &defaultConfigPath
		if err != nil {
			return nil, fmt.Errorf("failed to get config path: %v", err)
		}
	}

	fmt.Printf("config path: %s\n", *configPath)

	return &Options{
		Test:       *testFlag,
		Ntfy:       *ntfyFlag,
		Reset:      *resetFlag,
		Read:       *readFlag,
		DryRun:     *dryRunFlag,
		ConfigPath: *configPath,
	}, nil
}
