package battery

import "github.com/akryptic/battery-notifier/internal/logging"

func ReadLevel() (int, error) {
	logging.Trace("Reading battery level")
	level, err := batteryLevel()
	return int(level), err
}

func ReadStatus() (BatteryStatus, error) {
	logging.Trace("Reading battery charging status")
	status, err := isCharging()

	if err != nil {
		return "", err
	}

	if status {
		return Charging, nil
	}

	return Discharging, nil
}

func ReadBatteryState() (BatteryState, error) {
	logging.Trace("Reading complete battery state")
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
