package battery

import "errors"

type BatteryStatus = string

const (
	Charging    BatteryStatus = "Charging"
	Discharging BatteryStatus = "Discharging"
)

type BatteryState struct {
	Level  int
	Status BatteryStatus
}

var (
	// ErrNotAvailableAPI indicates that the current device/OS doesn't support such function.
	ErrNotAvailableAPI = errors.New("pref: not available api")

	// ErrNoSystemBattery indicates that the current device doesn't use batteries.
	//
	// Some APIs (like Android and JS) don't provide a mechanism to determine whether the machine uses batteries or not.
	// In such a case ErrNoSystemBattery will never be returned.
	ErrNoSystemBattery = errors.New("pref: device isn't battery-powered")
)
