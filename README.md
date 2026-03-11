# Projet 3A - Monitoring de la queue du Ru : WiFi Sniffer

## Installation

### Install Kismet

[documentation here](https://www.kismetwireless.net/docs/readme/installing/)

For Debian Trixie (compatible with the most recent versions of RaspberryPi OS) :

It is reccomended to install kismet as suid-root.

```bash
wget -O - https://www.kismetwireless.net/repos/kismet-release.gpg.key --quiet | gpg --dearmor | sudo tee /usr/share/keyrings/kismet-archive-keyring.gpg >/dev/null
echo 'deb [signed-by=/usr/share/keyrings/kismet-archive-keyring.gpg] https://www.kismetwireless.net/repos/apt/release/trixie trixie main' | sudo tee /etc/apt/sources.list.d/kismet.list >/dev/null
sudo apt update
sudo apt install kismet
```

If you installed kismet as suid-root you need to join the kismet group :

```bash
sudo usermod -aG kismet your-user-here
```

Create a service in `/etc/systemd/system/kismet.service` :

```INI
[Unit]
Description=Kismet Wireless Sniffer
After=network.target

[Service]
Type=simple
User=pi
Group=kismet
ExecStart=/usr/bin/kismet --no-line-wrap
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Configure Kismet

All config overrides are done in `/etc/kismet/kismet_site.conf`.

Create the log folder :

```bash
sudo mkdir -p /var/log/kismet
sudo chown your-user-here:kismet /var/log/kismet
sudo chmod 775 /var/log/kismet
```

Configure the logs folder :

```conf
log_prefix=/var/log/kismet/
```

Configure the network interface(s) used to capture traffic :

```conf
source=wlan1
```

Change the log timeout (logs will be kept this long) :

```conf
kis_log_device_timeout=1800
kis_log_packet_timeout=1800
kis_log_alert_timeout=1800
kis_log_message_timeout=1800
kis_log_snapshot_timeout=1800
```

Flag the logs as ephemeral :

```conf
kis_log_ephemeral_dangerous=true
```

### Start the Service

```bash
sudo systemctl daemon-reload
sudo systemctl enable kismet
sudo systemctl start kismet
```

Check the service :

```bash
sudo systemctl status kismet
```

## Usage

If you have go installed on the Raspberry Pi:

```bash
go run ./deviceCounter.go
```
