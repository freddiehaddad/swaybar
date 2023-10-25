package cpu

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/freddiehaddad/swaybar/pkg/descriptor"
	"github.com/freddiehaddad/swaybar/pkg/utils"
)

type CPU struct {
	Sensor   string
	Interval time.Duration
	Enabled  atomic.Bool
}

func New(sensor string, interval time.Duration) (*CPU, error) {
	cpu := &CPU{
		Sensor:   sensor,
		Interval: interval,
	}

	cpu.Update()
	return cpu, nil
}

func (c *CPU) Update() (descriptor.Descriptor, error) {
	log.Println("Updating", c.Sensor)
	descriptor := descriptor.Descriptor{
		Component: "cpu",
		Value:     "",
	}
	var sb strings.Builder

	sensor := fmt.Sprintf("/sys/class/hwmon/hwmon3/%s", c.Sensor)
	sensorValue, err := utils.GetSensorValue(sensor)
	if err != nil {
		log.Println("Error reading", sensor, err)
		return descriptor, err
	}

	tempCelcius:= utils.ReadSensorValue(sensorValue)
	sb.WriteString(fmt.Sprintf("CPU %5.1f °C", tempCelcius))
	descriptor.Value = sb.String()
	return descriptor, nil
}

func (c *CPU) Start(buffer chan descriptor.Descriptor) {
	c.Enabled.Store(true)

	go func() {
		for c.Enabled.Load() {
			descriptor, err := c.Update()
			if err != nil {
				log.Println("Error during update", err)
			} else {
				buffer <- descriptor
			}
			time.Sleep(c.Interval)
		}
	}()
}

func (c *CPU) Stop() {
	c.Enabled.Store(false)
}
