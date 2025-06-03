// go:build darwin

// The code is partially modified version of the code from the following repository.
// LINK: https://github.com/distatus/battery/blob/master/battery_darwin.go
// Copyright (C) 2016-2017,2023 Karol 'Kenji Takahashi' Woźniak
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package battery

import (
	"os/exec"

	plist "howett.net/plist"
)

type battery struct {
	Voltage           int
	CurrentCapacity   int `plist:"AppleRawCurrentCapacity"`
	MaxCapacity       int `plist:"AppleRawMaxCapacity"`
	DesignCapacity    int
	Amperage          int64
	FullyCharged      bool
	IsCharging        bool
	ExternalConnected bool
}

func readBatteries() ([]*battery, error) {
	out, err := exec.Command("ioreg", "-n", "AppleSmartBattery", "-r", "-a").Output()
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, nil
	}

	var data []*battery
	if _, err = plist.Unmarshal(out, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func batteryLevel() (uint8, error) {
	batts, err := readBatteries()
	if err != nil {
		return 0, err
	}

	if len(batts) == 0 {
		return 0, ErrNoSystemBattery
	}

	cap := uint8(float64(batts[0].CurrentCapacity) / float64(batts[0].MaxCapacity) * 100)

	return cap, nil
}

func isCharging() (bool, error) {

	batts, err := readBatteries()

	if err != nil {
		return false, err
	}

	if len(batts) == 0 {
		return false, ErrNoSystemBattery
	}

	return batts[0].IsCharging, nil
}
