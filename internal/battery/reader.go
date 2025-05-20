package battery

import (
	bat "gioui.org/x/pref/battery"
)

type BatteryStatus = string

const (
	Charging    BatteryStatus = "Charging"
	Discharging BatteryStatus = "Discharging"
)

type BatteryState struct {
	Level  int
	Status BatteryStatus
}

func ReadLevel() (int, error) {
	level, err := bat.Level()
	return int(level), err
}

func ReadStatus() (BatteryStatus, error) {
	status, err := bat.IsCharging()

	if err != nil {
		return "", err
	}

	if status {
		return Charging, nil
	}

	return Discharging, nil
}

func ReadBatteryState() (BatteryState, error) {
	level, err := ReadLevel()
	if err != nil {
		return BatteryState{}, err
	}

	status, err := ReadStatus()
	if err != nil {
		return BatteryState{}, err
	}

	return BatteryState{level, status}, nil
}
