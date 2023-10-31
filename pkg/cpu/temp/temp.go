package cputemp

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/freddiehaddad/swaybar/pkg/descriptor"
	"github.com/freddiehaddad/swaybar/pkg/utils"
)

type CPUTemp struct {
	Sensor   string
	Interval time.Duration
	Enabled  atomic.Bool
}

func New(sensor string, interval time.Duration) (*CPUTemp, error) {
	cpu := &CPUTemp{
		Sensor:   sensor,
		Interval: interval,
	}

	cpu.Update()
	return cpu, nil
}

func (c *CPUTemp) Update() (descriptor.Descriptor, error) {
	log.Printf("INFO: Updating CPU temperature sensor=%s", c.Sensor)
	descriptor := descriptor.Descriptor{
		Component: "cputemp",
		Value:     "",
	}
	var sb strings.Builder

	sensor := fmt.Sprintf("/sys/class/hwmon/hwmon3/%s", c.Sensor)
	sensorValue, err := utils.GetSensorValue(sensor)
	if err != nil {
		log.Printf("ERROR: GetSensorValue sensor=%s err=%s", sensor, err)
		return descriptor, err
	}

	tempCelcius:= utils.ReadSensorValue(sensorValue)
	sb.WriteString(fmt.Sprintf("CPU %5.1f °C", tempCelcius))
	descriptor.Value = sb.String()
	return descriptor, nil
}

func (c *CPUTemp) Start(buffer chan descriptor.Descriptor) {
	c.Enabled.Store(true)

	go func() {
		for c.Enabled.Load() {
			descriptor, err := c.Update()
			if err != nil {
				log.Printf("ERROR: Update err=%s", err)
			} else {
				buffer <- descriptor
			}
			time.Sleep(c.Interval)
		}
	}()
}

func (c *CPUTemp) Stop() {
	c.Enabled.Store(false)
}
