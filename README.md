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

- Time
- Network

## Building

```text
go build -o bin ./...
```

## Testing

```text
go test ./...
```

## Installation

After building the code, copy the generated binary to a directory of your choice
and update your sway config `bar` section to launch the code. In this example,
the program is copied to `$HOME/.config/sway/statusbar`.

```text
bar {
    status_command $HOME/.config/sway/statusbar
}
```
