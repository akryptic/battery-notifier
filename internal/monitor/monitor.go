package monitor

import (
	"fmt"
	"time"

	"github.com/akryptic/battery-notifier/internal/battery"
	"github.com/akryptic/battery-notifier/internal/config"
	"github.com/akryptic/battery-notifier/internal/logging"
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
	logging.Trace("Creating new monitor with config: LowBattery=%d, Critical=%d, Overcharge=%d",
		conf.LowBattery, conf.CriticalBattery, conf.OverchargeLimit)
	return &Monitor{
		Config:        conf,
		NotifiedState: &NotifiedState{false, false, false},
	}
}

func (m *Monitor) determineRunType() (RunType, battery.BatteryState, error) {
	batteryState, err := battery.ReadBatteryState()
	if err != nil {
		logging.Error("Failed to read battery state: %v", err)
		return NoRun, batteryState, err
	}

	logging.Debug("Battery: %d%% %s", batteryState.Level, batteryState.Status)

	if batteryState.Level <= m.Config.CriticalBattery && batteryState.Status != battery.Charging {
		logging.Debug("CRITICAL: Battery critically low (%d%% <= %d%%)", batteryState.Level, m.Config.CriticalBattery)
		return CriticalBatteryRun, batteryState, nil
	}

	if batteryState.Level <= m.Config.LowBattery && batteryState.Status != battery.Charging {
		logging.Debug("LOW: Battery low (%d%% <= %d%%)", batteryState.Level, m.Config.LowBattery)
		return LowBatteryRun, batteryState, nil
	}

	if batteryState.Level >= m.Config.OverchargeLimit && batteryState.Status == battery.Charging {
		logging.Debug("OVERCHARGING: Battery overcharging (%d%% >= %d%%)", batteryState.Level, m.Config.OverchargeLimit)
		return OverchargeLimitRun, batteryState, nil
	}

	return NoRun, batteryState, nil
}

// handles notification state and sending notifications
func (m *Monitor) ProcessNotifications() (RunType, error) {

	logging.Trace("Processing notifications")

	runType, batteryState, err := m.determineRunType()

	logging.Trace("Run: %s", runType)
	logging.Trace("Battery state: %+v", batteryState)
	logging.Trace("Notified state: %+v", m.NotifiedState)

	if err != nil {
		logging.Error("Failed to determine run type: %v", err)
		return NoRun, err
	}

	if runType == NoRun {
		logging.Trace("No notifications needed")
		return NoRun, nil
	}

	if runType == LowBatteryRun {
		if !m.NotifiedState.Low {
			logging.Debug("Sending low battery notification (%d%%)", batteryState.Level)
			err := notification.SendNotification(
				"Battery Low",
				fmt.Sprintf("%d%% remaining. Please plug in.", batteryState.Level),
				m.Config,
			)
			if err != nil {
				logging.Error("Failed to send low battery notification: %v", err)
				return NoRun, err
			}

			if m.Config.EnableSound {
				err := sound.Play("low", m.Config)
				if err != nil {
					logging.Warn("Failed to play low battery sound: %v", err)
				}
			}

			m.NotifiedState.Low = true
			return LowBatteryRun, nil
		} else {
			logging.Debug("Low Battery notification already sent, skipping")
		}
	}

	if batteryState.Level > m.Config.LowBattery && m.NotifiedState.Low {
		logging.Debug("Battery level recovered above low threshold (%d%% > %d%%), resetting low notification state",
			batteryState.Level, m.Config.LowBattery)
		m.NotifiedState.Low = false
		return NoRun, nil
	}

	if runType == CriticalBatteryRun {
		if !m.NotifiedState.Critical {
			logging.Debug("Sending critical battery notification (%d%%)", batteryState.Level)
			err := notification.SendNotification(
				"Battery Critically Low",
				fmt.Sprintf("%d%% remaining! System may shut down.", batteryState.Level),
				m.Config,
			)
			if err != nil {
				logging.Error("Failed to send critical battery notification: %v", err)
				return NoRun, err
			}

			if m.Config.EnableSound {
				err := sound.Play("low", m.Config)
				if err != nil {
					logging.Warn("Failed to play critical battery sound: %v", err)
					// continue even if sound fails
				}
			}

			m.NotifiedState.Critical = true
			return CriticalBatteryRun, nil
		} else {

			logging.Debug("Critical Battery notification already sent, skipping")
		}
	}

	if batteryState.Level > m.Config.CriticalBattery && m.NotifiedState.Critical {
		logging.Debug("Battery level recovered above critical threshold (%d%% > %d%%), resetting critical notification state",
			batteryState.Level, m.Config.CriticalBattery)
		m.NotifiedState.Critical = false
		return NoRun, nil
	}

	if runType == OverchargeLimitRun {
		if !m.NotifiedState.Overcharge {
			logging.Debug("Sending overcharge notification (%d%%)", batteryState.Level)
			err := notification.SendNotification(
				"Battery Overcharging",
				fmt.Sprintf(
					"%d%% charged. Consider unplugging to preserve battery health.",
					batteryState.Level,
				),
				m.Config,
			)
			if err != nil {
				logging.Error("Failed to send overcharge notification: %v", err)
				return NoRun, nil
			}

			if m.Config.EnableSound {
				err := sound.Play("overcharge", m.Config)
				if err != nil {
					logging.Warn("Failed to play overcharge sound: %v", err)
					// continue even if sound fails
				}
			}

			m.NotifiedState.Overcharge = true
			return OverchargeLimitRun, nil
		} else {
			logging.Debug("Overcharge notification already sent, skipping")
		}
	}

	if batteryState.Level < m.Config.OverchargeLimit && m.NotifiedState.Overcharge {
		logging.Debug("Battery level dropped below overcharge threshold (%d%% < %d%%), resetting overcharge notification state",
			batteryState.Level, m.Config.OverchargeLimit)
		m.NotifiedState.Overcharge = false
		return NoRun, nil
	}

	return NoRun, nil
}

// continuous monitoring of battery state
func (m *Monitor) StartMonitoring() {
	logging.Debug("Starting battery monitoring with %d second intervals", m.Config.CheckInterval)

	for {
		logging.Info("Running battery check")
		runType, err := m.ProcessNotifications()
		if err != nil {
			logging.Error("Monitoring failed: %v", err)
			return
		}
		logging.Info("Run: %s", runType)

		logging.Trace("Sleeping for %d seconds", m.Config.CheckInterval)
		time.Sleep(time.Duration(m.Config.CheckInterval) * time.Second)
	}
}

// dry-run
func (m *Monitor) RunOnce() (RunType, error) {
	logging.Info("Running battery check (dry-run mode)")
	runType, err := m.ProcessNotifications()
	if err != nil {
		logging.Error("Dry-run failed: %v", err)
	}
	return runType, err
}
