mod config;
use config::Config;
use notify_rust::Notification;
use std::{env, fs, path::PathBuf, process::Command, thread, time::Duration};

fn main() {
    let config_path = PathBuf::from(format!(
        "{}/.config/battery-notifier/config.toml",
        std::env::var("HOME").unwrap()
    ));

    let args: Vec<String> = env::args().collect();
    let config = Config::load(&config_path);

    // If "--test" flag is passed, show a test notification and exit
    if args.len() > 1 {
        if args.contains(&"--test".to_string()) {
            send_notification(
                "Battery Notifier Test",
                "This is a test notification.",
                &config,
            );
        }

        if args.contains(&"--notify".to_string()) {
            if let Some(topic) = &config.ntfy_topic {
                send_noti_via_ntfy(
                    "This is a test ntfy notification.",
                    topic,
                    &config.ntfy_server,
                );
            } else {
                eprintln!("No ntfy topic set in config.");
            }
        }

        if args.contains(&"--reset".to_string()) {
            Config::generate_default_config(&config_path);
            print!("Default config generated at: {:?}\n", config_path);
        }
        return;
    }

    let mut notified_low = false;
    let mut notified_critical = false;
    let mut notified_overcharge = false;

    loop {
        if let Ok(capacity) = read_battery_capacity() {
            let status = read_battery_status().unwrap_or_else(|| "Unknown".to_string());

            if capacity <= config.low_battery && !notified_low && status != "Charging" {
                send_notification(
                    "Battery Low",
                    &format!("{}% remaining. Please plug in.", capacity),
                    &config,
                );
                if let Some(topic) = &config.ntfy_topic {
                    send_noti_via_ntfy(
                        &format!("Battery Low: {}% remaining.", capacity),
                        topic,
                        &config.ntfy_server,
                    );
                }
                notified_low = true;
            }

            if capacity > config.low_battery && notified_low {
                notified_low = false;
            }

            if capacity <= config.critical_battery && !notified_critical && status != "Charging" {
                send_notification(
                    "Battery Critically Low",
                    &format!("{}% remaining! System may shut down.", capacity),
                    &config,
                );
                notified_critical = true;
            }

            if capacity > config.critical_battery && notified_critical {
                notified_critical = false;
            }

            if capacity >= config.overcharge_limit && status == "Charging" && !notified_overcharge {
                send_notification(
                    "Battery Overcharging",
                    &format!(
                        "{}% charged. Consider unplugging to preserve battery health.",
                        capacity
                    ),
                    &config,
                );
                if let Some(topic) = &config.ntfy_topic {
                    send_noti_via_ntfy(
                        &format!("Battery Overcharging: {}% charged.", capacity),
                        topic,
                        &config.ntfy_server,
                    );
                }
                notified_overcharge = true;
            }
            if capacity < config.overcharge_limit && notified_overcharge {
                notified_overcharge = false;
            }
        }

        thread::sleep(Duration::from_secs(config.check_interval));
    }
}

fn read_battery_capacity() -> Result<u8, std::io::Error> {
    let content = fs::read_to_string("/sys/class/power_supply/BAT0/capacity")?;
    content.trim().parse().map_err(|_| {
        std::io::Error::new(std::io::ErrorKind::InvalidData, "Failed to parse capacity")
    })
}

fn read_battery_status() -> Option<String> {
    fs::read_to_string("/sys/class/power_supply/BAT0/status")
        .ok()
        .map(|s| s.trim().to_string())
}

fn send_notification(title: &str, message: &str, config: &Config) {
    let _ = Notification::new()
        .summary(title)
        .body(message)
        .timeout(5000)
        .show();

    // If ntfy.sh is configured, send a notification

    if config.enable_sound {
        let sound_path = config.sound_file.as_deref().unwrap_or("bell");
        let volume = format!("--volume={}", config.sound_volume.clamp(0, 100));

        let result = Command::new("canberra-gtk-play")
            .arg("-f")
            .arg(sound_path)
            .arg(volume)
            .spawn();

        match result {
            Ok(mut child) => {
                // Wait for the sound process to finish and check for errors
                if let Err(e) = child.wait() {
                    eprintln!("Failed to execute sound command: {}", e);
                    std::process::exit(1); // Exit if sound command fails
                }
            }
            Err(e) => {
                eprintln!("Failed to start sound command: {}", e);
                std::process::exit(1); // Exit if starting the command fails
            }
        }
    }
}

fn send_noti_via_ntfy(message: &str, topic: &str, server: &str) {
    if !topic.is_empty() {
        let ntfy_url = format!("{}/{}", server, topic);
        let result = Command::new("curl")
            .arg("-d")
            .arg(format!("{}", message))
            .arg(ntfy_url)
            .output();

        if let Err(e) = result {
            eprintln!("Failed to send ntfy notification: {}", e);
        }
    }
}
