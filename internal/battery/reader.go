package battery

import (
	"os"
	"strconv"
	"strings"
)

func ReadCapacity() (int, error) {
	content, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		return 0, err
	}

	capacity, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil {
		return 0, err
	}

	return capacity, nil
}

func ReadStatus() (string, error) {
	content, err := os.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {
		return "", err
	}

	status := strings.TrimSpace(string(content))
	return status, nil
}
