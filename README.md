# Battery Notifier

Battery Notifier is a simple GoLang-based utility that monitors battery levels and sends notifications when the battery is low, critically low, or overcharged. It supports local notifications and optional push notifications via [ntfy.sh](https://ntfy.sh/).

## Features

- Notifies when the battery is **low** (default: 20%)
- Notifies when the battery is **critically low** (default: 10%)
- Notifies when the battery is **overcharged** while plugged in (default: 80%)
- Plays a sound notification (configurable)
- Supports push notifications via [ntfy.sh](https://ntfy.sh/)
- Configurable check intervals
- Automatically generates a default configuration file
- Cross-platform support (Linux and Windows)

## Installation

### Prerequisites

#### Linux

- GoLang 1.24.1 or later (for building from source)
- `libnotify` (for local notifications)
- Audio libraries for sound playback (requirements for beep/oto packages)

#### Windows

- GoLang 1.24.1 or later (for building from source)

### Building from Source

```sh
# Clone the repository
git clone https://github.com/akryptic/battery-notifier.git
cd battery-notifier

# Build the application
go build -o battery-notifier

# Linux: Move the binary to a directory in your PATH
sudo mv battery-notifier /usr/local/bin/

# Windows: The build produces battery-notifier.exe
# Move it to your preferred location
```

### Reducing binary size

To reduce the binary size, you can use `UPX`. UPX is a free, portable, extendable, high-performance executable packer. For more details on how to use UPX, please refer to its [GitHub repository](https://github.com/upx/upx).

## Configuration

Battery Notifier uses a configuration file located at:

### Linux

```
~/.config/battery-notifier/config.toml
```

### Windows

```
%APPDATA%\battery-notifier\config.toml
```

If this file does not exist, the application will generate a default configuration.

### Default Config (`config.toml`):

```toml
# Battery threshold settings
low_battery = 20
critical_battery = 10
overcharge_limit = 80

# Notification settings
enable_sound = true
low_sound_file = ""         # Optional: path to custom sound for low battery (leave empty to use default)
overcharge_sound_file = ""  # Optional: path to custom sound for overcharge alert (leave empty to use default)
sound_volume = 80           # Volume level (0-100)

# Interval settings
check_interval = 60         # Time in seconds between battery status checks (min: 1, max: 300)

# ntfy.sh settings (leave empty if not using ntfy)
ntfy_topic = ""
ntfy_server = "https://ntfy.sh"
ntfy_access_token = ""      # Optional: for private topics
```

## Usage

### Running Manually

#### Linux and Windows

```sh
battery-notifier
```

### Running on Startup

#### Linux (Hyprland Example)

Add this line to your `hyprland.conf`:

```ini
exec-once = battery-notifier &
```

#### Linux (Systemd Example)

Create a new service file at `/etc/systemd/system/battery-notifier.service`:

```ini
[Unit]
Description=Battery Notifier Service
After=network.target

[Service]
ExecStart=/usr/local/bin/battery-notifier
Restart=on-failure
User=YOUR_USERNAME

[Install]
WantedBy=multi-user.target
```

Replace `YOUR_USERNAME` with your actual username.

Reload the systemd daemon and enable the service:

```sh
systemctl daemon-reload
systemctl enable battery-notifier.service
systemctl start battery-notifier.service
```

Check the status of the service:

```sh
systemctl status battery-notifier.service
```

#### Windows (Task Scheduler)

1. Open Task Scheduler (search for "Task Scheduler" in the Start menu)
2. Click "Create Basic Task..."
3. Enter a name (e.g., "Battery Notifier") and description
4. Select "When I log on" as the trigger
5. Select "Start a program" as the action
6. Browse to the location of your battery-notifier.exe
7. Click "Finish"

#### Windows (Alternative Startup Methods)

**Method 1: Startup Folder**

1. Press `Win + R` to open the Run dialog
2. Type `shell:startup` and press Enter
3. Create a shortcut to your battery-notifier.exe in this folder

**Method 2: Registry**

1. Press `Win + R` to open the Run dialog
2. Type `regedit` and press Enter
3. Navigate to `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run`
4. Right-click in the right pane and select New > String Value
5. Name it "Battery Notifier"
6. Set the value to the full path to your battery-notifier.exe

**Method 3: Batch Script + Shortcut**

1. Create a batch file (e.g., `start_battery_notifier.bat`) with the following content:
   ```batch
   @echo off
   start "" "C:\path\to\battery-notifier.exe"
   ```
2. Create a shortcut to this batch file in your startup folder (`shell:startup`)
3. Set the shortcut to run minimized

### Testing Notifications

```sh
# Local notification test (with sound if enabled)
battery-notifier --test

# Print current battery status
battery-notifier --read

# Send a test notification via ntfy.sh (requires ntfy_topic in config)
battery-notifier --ntfy

# Reset config to default values
# USE CAREFULLY
battery-notifier --reset

# Execute a single check without continuous monitoring
battery-notifier --dry-run
```

## License

This project is licensed under the MIT License.

## Contributions

Contributions and suggestions are welcome! Feel free to open an issue or submit a pull request.
