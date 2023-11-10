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

const (
	Name  = "coretemp"
	Label = "package"
)

type CPUTemp struct {
	Name       string
	Label      string
	SensorPath string
	Interval   time.Duration
	Enabled    atomic.Bool
}

func (c *CPUTemp) init() error {
	return nil
}

// Attempt to create a sensor with data coming from the hwmon (name) and
// sensor (label). If either string is empty, internal default values will
// be assumed.
func New(name, label string, interval time.Duration) (*CPUTemp, error) {
	cpu := &CPUTemp{
		Interval: interval,
	}

	if len(name) == 0 {
		cpu.Name = Name
	} else {
		cpu.Name = name
	}

	if len(label) == 0 {
		cpu.Label = Label
	} else {
		cpu.Label = label
	}

	sensorPath, err := utils.FindSensorPath(cpu.Name, cpu.Label)
	if err != nil {
		return nil, err
	}

	cpu.SensorPath = sensorPath
	cpu.Update()
	return cpu, nil
}

func (c *CPUTemp) Update() (descriptor.Descriptor, error) {
	log.Printf("INFO: Updating CPU temperature sensor=%s", c.Label)
	descriptor := descriptor.Descriptor{
		Component: "cputemp",
		Value:     "",
	}
	var sb strings.Builder

	sensorValue, err := utils.GetSensorValue(c.SensorPath)
	if err != nil {
		log.Printf("ERROR: GetSensorValue sensor=%s err=%s", c.SensorPath, err)
		return descriptor, err
	}

	tempCelcius := utils.ReadSensorValue(sensorValue)
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
