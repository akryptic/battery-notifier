# Battery Notifier

Battery Notifier is a simple Rust-based utility that monitors battery levels and sends notifications when the battery is low, critically low, or overcharged. It supports local notifications and optional push notifications via [ntfy.sh](https://ntfy.sh/).

## Features

- Notifies when the battery is **low** (default: 20%)
- Notifies when the battery is **critically low** (default: 10%)
- Notifies when the battery is **overcharged** while plugged in (default: 80%)
- Plays a sound notification (configurable)
- Supports push notifications via [ntfy.sh](https://ntfy.sh/)
- Configurable check intervals
- Automatically generates a default configuration file

## Installation

### Prerequisites

- Rust (for building from source)
- `canberra-gtk-play` (for playing notification sounds)
- `libnotify` (for local notifications)
- `curl` (for ntfy.sh support)

### Building from Source

```sh
# Clone the repository
git clone https://github.com/ardentkilnfire/battery-notifier.git
cd battery-notifier

# Build the application
cargo build --release

# Move the binary to a directory in your PATH
sudo mv target/release/battery-notifier /usr/local/bin/
```

### Reducing binary size

To reduce the binary size, you can use `UPX`. UPX is a free, portable, extendable, high-performance executable packer. For more details on how to use UPX, please refer to its [GitHub repository](https://github.com/upx/upx).

## Configuration

Battery Notifier uses a configuration file located at:

```
~/.config/battery-notifier/config.toml
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
sound_file = "/usr/share/sounds/freedesktop/stereo/bell.oga"
sound_volume = 100  # Volume level (0-100)

# Interval settings
check_interval = 60  # Time in seconds between battery status checks (min: 1, max: 300)

# ntfy.sh settings (leave empty if not using ntfy)
ntfy_topic = ""
ntfy_server = "https://ntfy.sh"
```

## Usage

### Running Manually

```sh
battery-notifier
```

### Running on Startup (Hyprland Example)

Add this line to your `hyprland.conf`:

```ini
exec-once = battery-notifier &
```

### Running on Startup (Systemd Example)

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

### Testing Notifications

- Local notification test:

```sh
battery-notifier --test
```

- ntfy.sh test (requires `ntfy_topic` to be set in the config):

```sh
battery-notifier --notify
```

- Reset config to default:

```sh
battery-notifier --reset
```

## License

This project is licensed under the MIT License.

## Contributions

Contributions and suggestions are welcome! Feel free to open an issue or submit a pull request.
