package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	LowBattery      int `toml:"low_battery"`
	CriticalBattery int `toml:"critical_battery"`
	OverchargeLimit int `toml:"overcharge_limit"`

	EnableSound         bool   `toml:"enable_sound"`
	LowSoundFile        string `toml:"low_sound_file"`
	OverchargeSoundFile string `toml:"overcharge_sound_file"`
	SoundVolume         int    `toml:"sound_volume"`

	CheckInterval int `toml:"check_interval"`

	NtfyTopic       string `toml:"ntfy_topic"`
	NtfyServer      string `toml:"ntfy_server"`
	NtfyAccessToken string `toml:"ntfy_access_token"`
}

// utility
func (c *Config) GetNtfyUrl() string {
	return fmt.Sprintf("%s/%s", c.NtfyServer, c.NtfyTopic)
}

func (c *Config) Validate() string {

	if c.LowBattery < 0 {
		return "Low battery threshold cannot be negative"
	}

	if c.CriticalBattery < 0 {
		return "Critical battery threshold cannot be negative"
	}

	if c.OverchargeLimit < 0 {
		return "Overcharge limit cannot be negative"
	}

	if c.EnableSound && c.SoundVolume < 0 {
		return "Sound volume cannot be negative"
	}

	if c.EnableSound && c.SoundVolume > 100 {
		return "Sound volume cannot be greater than 100"
	}

	if c.CheckInterval < 1 {
		return "Check interval cannot be less than 1 second"
	}

	if c.CheckInterval > 300 {
		return "Check interval cannot be greater than 300 seconds"
	}

	if c.SoundVolume < 0 || c.SoundVolume > 100 {
		return "Sound volume must be between 0 and 100"
	}

	return ""
}

func Load(path string) (config Config, err error) {

	// if config file doesn't exist, generate a default config and return
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config := GenerateDefaultConfig(path)
		return config, nil
	}

	_, err = toml.DecodeFile(path, &config)
	return config, err
}

func GetDefaultConfig() Config {
	return Config{
		LowBattery:      20,
		CriticalBattery: 10,
		OverchargeLimit: 80,

		EnableSound:         true,
		LowSoundFile:        "",
		OverchargeSoundFile: "",
		SoundVolume:         80,

		CheckInterval: 60,

		NtfyTopic:       "",
		NtfyServer:      "https://ntfy.sh",
		NtfyAccessToken: "",
	}
}

func GenerateDefaultConfig(path string) Config {
	config := GetDefaultConfig()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = toml.NewEncoder(file).Encode(config)
	if err != nil {
		panic(err)
	}

	return config
}

// cross-platform default path for the config file
func GetDefaultConfigPath() (string, error) {
	userConfigDir, err := os.UserConfigDir()

	if err != nil {
		// fallback to HOME
		userConfigDir = os.Getenv("HOME")
		if userConfigDir == "" {
			return "", fmt.Errorf("couldn't determine config directory: %v", err)
		}
		userConfigDir = filepath.Join(userConfigDir, ".config")
	}

	return filepath.Join(userConfigDir, "battery-notifier", "config.toml"), nil
}
