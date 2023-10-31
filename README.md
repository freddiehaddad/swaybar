# Swaybar

A highly parallelized swaybar app written in Go.

The modules (or components) are each responsible for populating the status bar
with specific data. For example, the network component provides data throughput
rates for a given interface. Similarly, the time component returns time and date
values.

Each component runs in its own thread and reports updates on a specified
interval. For example, the network component can be set to publish updates every
second, millisecond, ten microseconds, etc.

Whenever an update is pushed, the status bar will get updated allowing each
components report updates at different intervals.

List of components currently implemented included:

- CPU Temperature
- CPU Utilization
- GPU Temperature
- Time
- Network Throughput

## Sample Output

```text
D  977.58 Kbps U   36.79 Mbps | CPU   0.9% | CPU  26.0 °C | GPU  58.0 °C | Wednesday, 25-Oct-23 16:14:20 PDT
```

## Building

```text
go build -v -o bin ./...
```

## Testing

```text
go test -v ./...
```

## Installation

After building the code, copy the generated binary to a directory of your choice
and update your sway config `bar` section to launch the code. In this example,
the program is copied to `$HOME/.config/sway/statusbar`.

The program will look for a `config.yml` file in a config directory in the same
location as the executable.

## Configuration

Refer to the `config/config.yml` reference file for configuring the status bar.

```text
bar {
    status_command $HOME/.config/sway/statusbar
}
```
