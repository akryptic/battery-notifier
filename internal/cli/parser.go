package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/akamensky/argparse"
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/logging"
)

type Options struct {
	Test         bool
	Ntfy         bool
	Reset        bool
	Read         bool
	DryRun       bool
	ConfigPath   string
	VerboseLevel int
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

	verboseLevel := parser.Int("v", "verbose", &argparse.Options{
		Required: false,
		Help:     "Verbose logging level (0=quiet, 1=info, 2=debug, 3=trace)",
		Default:  1,
		Validate: func(args []string) error {
			if len(args) > 0 {
				levelString := args[0]
				level, err := strconv.Atoi(levelString)
				if err != nil {
					return fmt.Errorf("verbose level must be an integer")
				}
				if level < 0 || level > 3 {
					return fmt.Errorf("verbose level must be between 0 and 3")
				}
			}
			return nil
		},
	})

	err := parser.Parse(os.Args)
	if err != nil {
		return nil, err
	}

	logging.SetLevel(*verboseLevel)

	// if config path was not passed, use the default path
	if *configPath == "" {
		defaultConfigPath, err := config.GetDefaultConfigPath()

		configPath = &defaultConfigPath
		if err != nil {
			return nil, fmt.Errorf("failed to get config path: %v", err)
		}
	}

	return &Options{
		Test:         *testFlag,
		Ntfy:         *ntfyFlag,
		Reset:        *resetFlag,
		Read:         *readFlag,
		DryRun:       *dryRunFlag,
		ConfigPath:   *configPath,
		VerboseLevel: *verboseLevel,
	}, nil
}
