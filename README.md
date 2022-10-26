# swaybar.sh

This script populates the swaybar with information about your system using the
JSON protocol.

Features:

- Date
- System uptime
- CPU temperature
- Memory
- Network IP
- Network bandwith

To use, add or update the `status_command` in swaybar section of your sway
config.

NOTE: The script might require updating the network device name and path to the
temperature data.

```text
bar {
    status_command $HOME/.config/sway/scripts/swaybar.sh
}
```

Sample output:

```text
Rx: 6.31 Tx: 984.86 (Mbit/s) | E: 192.168.1.150/24 1000 (Mbit/s) |\
T: 31 Gb U: 2 Gb F: 24 Gb | CPU: 29.0 C | U: 9:30 | Wed Oct 26 16:34
```
