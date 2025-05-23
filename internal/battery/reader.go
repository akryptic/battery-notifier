package battery

func ReadLevel() (int, error) {
	level, err := batteryLevel()
	return int(level), err
}

func ReadStatus() (BatteryStatus, error) {
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
