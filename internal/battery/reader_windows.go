// go:build windows

package battery

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	_Kernel32 = windows.NewLazySystemDLL("kernel32")

	_GetSystemPowerStatus = _Kernel32.NewProc("GetSystemPowerStatus")

	_BatteryFlagCharging        byte = 8
	_BatteryFlagNoSystemBattery byte = 128
	_BatteryFlagUnknownStatus   byte = 255

	_BatteryLifePercentUnknown byte = 255
)

func batteryLevel() (uint8, error) {
	resp, err := powerStatus()
	if err != nil {
		return 100, ErrNotAvailableAPI
	}

	if resp.BatteryLifePercent == _BatteryLifePercentUnknown || resp.BatteryFlag&_BatteryFlagNoSystemBattery > 0 {
		return 100, ErrNoSystemBattery
	}

	return resp.BatteryLifePercent, nil
}

func isCharging() (bool, error) {
	resp, err := powerStatus()
	if err != nil || resp.BatteryFlag == _BatteryFlagUnknownStatus {
		return false, ErrNotAvailableAPI
	}

	if resp.BatteryFlag&_BatteryFlagNoSystemBattery > 0 {
		return true, ErrNoSystemBattery
	}

	return resp.BatteryFlag&_BatteryFlagCharging > 0, nil
}

type _SystemPowerStatus struct {
	ACLineStatus        byte
	BatteryFlag         byte
	BatteryLifePercent  byte
	SystemStatusFlag    byte
	BatteryLifeTime     int32
	BatteryFullLifeTime int32
}

func powerStatus() (resp _SystemPowerStatus, err error) {
	r, _, err := _GetSystemPowerStatus.Call(uintptr(unsafe.Pointer(&resp)))
	if r == 0 {
		return resp, err
	}

	return resp, nil
}
