use serde::Deserialize;
use std::{fs, path::PathBuf};

#[derive(Debug, Deserialize)]
pub struct Config {
    pub low_battery: u8,
    pub critical_battery: u8,
    pub overcharge_limit: u8,
    pub enable_sound: bool,
    pub sound_file: Option<String>,
    pub sound_volume: u8,
    pub check_interval: u64,
    pub ntfy_topic: Option<String>,
}

impl Config {
    pub fn load(config_path: &PathBuf) -> Self {
        if !config_path.exists() {
            println!("Config file not found, creating a default one.");
            Self::generate_default_config(&config_path);
        }

        if let Ok(contents) = fs::read_to_string(&config_path) {
            if let Ok(mut config) = toml::from_str::<Config>(&contents) {
                // Clamp the check interval to a reasonable range
                config.check_interval = config.check_interval.clamp(1, 300);
                return config;
            } else {
                eprintln!("Failed to parse config file, using default config.");
            }
        }

        // Default config (if file missing or invalid)
        Config {
            low_battery: 20,
            critical_battery: 10,
            overcharge_limit: 80,
            enable_sound: true,
            sound_file: Some("/usr/share/sounds/freedesktop/stereo/bell.oga".to_string()),
            sound_volume: 100,
            check_interval: 60,
            ntfy_topic: None,
        }
    }

    pub fn generate_default_config(path: &PathBuf) {
        let default_config = r#"
# Battery threshold settings
    low_battery = 20
    critical_battery = 10
    overcharge_limit = 80

# Notification settings
    enable_sound = true
    sound_file = "/usr/share/sounds/freedesktop/stereo/bell.oga"
    sound_volume = 100  # Volume level (0-100)

# Interval settings
    check_interval = 60  # Time in seconds between battery status checks (min: 1, max: 300)

# ntfy.sh settings (leave empty if not using ntfy)
    ntfy_topic = ""
"#;

        if let Err(e) = fs::create_dir_all(path.parent().unwrap()) {
            eprintln!("Failed to create config directory: {}", e);
            return;
        }

        // Write the default config to the file
        if let Err(e) = fs::write(path, default_config) {
            eprintln!("Failed to create config file: {}", e);
        }
    }
}
