package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/akryptic/battery-notifier/internal/battery"
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/notification"
	"github.com/akryptic/battery-notifier/internal/sound"
)

// represents the type of notification run
type RunType string

const (
	NoRun              RunType = "No Run"
	LowBatteryRun      RunType = "Low Battery Run"
	CriticalBatteryRun RunType = "Critical Battery Run"
	OverchargeLimitRun RunType = "Overcharging Limit Run"
)

type NotifiedState struct {
	Low        bool
	Critical   bool
	Overcharge bool
}

type Monitor struct {
	Config        *config.Config
	NotifiedState *NotifiedState
}

func NewMonitor(conf *config.Config) *Monitor {
	return &Monitor{
		Config:        conf,
		NotifiedState: &NotifiedState{false, false, false},
	}
}

func (m *Monitor) determineRunType() (RunType, battery.BatteryState, error) {
	batteryState, err := battery.ReadBatteryState()
	if err != nil {
		log.Printf("ERROR: Failed to read battery state: %v", err)
		return NoRun, batteryState, err
	}

	fmt.Printf("[%s] Battery: %d%% %s\n", time.Now().Format("2006-01-02 15:04:05"), batteryState.Level, batteryState.Status)

	if batteryState.Level <= m.Config.CriticalBattery && batteryState.Status != battery.Charging {
		log.Printf("CRITICAL: Battery critically low (%d%% <= %d%%)", batteryState.Level, m.Config.CriticalBattery)
		return CriticalBatteryRun, batteryState, nil
	}

	if batteryState.Level <= m.Config.LowBattery && batteryState.Status != battery.Charging {
		log.Printf("WARNING: Battery low (%d%% <= %d%%)", batteryState.Level, m.Config.LowBattery)
		return LowBatteryRun, batteryState, nil
	}

	if batteryState.Level >= m.Config.OverchargeLimit && batteryState.Status == battery.Charging {
		log.Printf("WARNING: Battery overcharging (%d%% >= %d%%)", batteryState.Level, m.Config.OverchargeLimit)
		return OverchargeLimitRun, batteryState, nil
	}

	return NoRun, batteryState, nil
}

// handles notification state and sending notifications
func (m *Monitor) ProcessNotifications() (RunType, error) {
	runType, batteryState, err := m.determineRunType()

	if err != nil {
		log.Printf("ERROR: Failed to determine run type: %v", err)
		return NoRun, err
	}

	if runType == NoRun {
		return NoRun, nil
	}

	if runType == LowBatteryRun && !m.NotifiedState.Low {
		log.Printf("Sending low battery notification (%d%%)", batteryState.Level)
		err := notification.SendNotification(
			"Battery Low",
			fmt.Sprintf("%d%% remaining. Please plug in.", batteryState.Level),
			m.Config,
		)
		if err != nil {
			log.Printf("ERROR: Failed to send low battery notification: %v", err)
			return NoRun, err
		}

		if m.Config.EnableSound {
			err := sound.Play("low", m.Config)
			if err != nil {
				log.Printf("WARNING: Failed to play low battery sound: %v", err)
				// continue even if sound fails
			}
		}

		m.NotifiedState.Low = true
		return LowBatteryRun, nil
	}

	if batteryState.Level > m.Config.LowBattery && m.NotifiedState.Low {
		m.NotifiedState.Low = false
		return NoRun, nil
	}

	if runType == CriticalBatteryRun && !m.NotifiedState.Critical {
		log.Printf("Sending critical battery notification (%d%%)", batteryState.Level)
		err := notification.SendNotification(
			"Battery Critically Low",
			fmt.Sprintf("%d%% remaining! System may shut down.", batteryState.Level),
			m.Config,
		)
		if err != nil {
			log.Printf("ERROR: Failed to send critical battery notification: %v", err)
			return NoRun, err
		}

		if m.Config.EnableSound {
			err := sound.Play("low", m.Config)
			if err != nil {
				log.Printf("WARNING: Failed to play critical battery sound: %v", err)
				// continue even if sound fails
			}
		}

		m.NotifiedState.Critical = true
		return CriticalBatteryRun, nil
	}

	if batteryState.Level > m.Config.CriticalBattery && m.NotifiedState.Critical {
		m.NotifiedState.Critical = false
		return NoRun, nil
	}

	if runType == OverchargeLimitRun && !m.NotifiedState.Overcharge {
		log.Printf("Sending overcharge notification (%d%%)", batteryState.Level)
		err := notification.SendNotification(
			"Battery Overcharging",
			fmt.Sprintf(
				"%d%% charged. Consider unplugging to preserve battery health.",
				batteryState.Level,
			),
			m.Config,
		)
		if err != nil {
			log.Printf("ERROR: Failed to send overcharge notification: %v", err)
			return NoRun, nil
		}

		if m.Config.EnableSound {
			err := sound.Play("overcharge", m.Config)
			if err != nil {
				log.Printf("WARNING: Failed to play overcharge sound: %v", err)
				// continue even if sound fails
			}
		}

		m.NotifiedState.Overcharge = true
		return OverchargeLimitRun, nil
	}

	if batteryState.Level < m.Config.OverchargeLimit && m.NotifiedState.Overcharge {
		m.NotifiedState.Overcharge = false
		return NoRun, nil
	}

	return NoRun, nil
}

// continuous monitoring of battery state
func (m *Monitor) StartMonitoring() {
	log.Printf("Starting battery monitoring with %d second intervals", m.Config.CheckInterval)

	for {
		runType, err := m.ProcessNotifications()
		if err != nil {
			log.Printf("ERROR: Monitoring failed: %v", err)
			fmt.Println(err)
			return
		}

		fmt.Printf("Run: %s\n", runType)
		time.Sleep(time.Duration(m.Config.CheckInterval) * time.Second)
	}
}

// dry-run
func (m *Monitor) RunOnce() (RunType, error) {
	log.Println("Running battery check (dry-run mode)")
	runType, err := m.ProcessNotifications()
	if err != nil {
		log.Printf("ERROR: Dry-run failed: %v", err)
	}
	return runType, err
}
