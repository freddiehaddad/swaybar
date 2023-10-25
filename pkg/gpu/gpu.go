package gpu

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/freddiehaddad/swaybar/pkg/descriptor"
	"github.com/freddiehaddad/swaybar/pkg/utils"
)

type GPU struct {
	Sensor   string
	Interval time.Duration
	Enabled  atomic.Bool
}

func New(sensor string, interval time.Duration) (*GPU, error) {
	gpu := &GPU{
		Sensor:   sensor,
		Interval: interval,
	}

	gpu.Update()
	return gpu, nil
}

func (c *GPU) Update() (descriptor.Descriptor, error) {
	log.Println("Updating", c.Sensor)
	descriptor := descriptor.Descriptor{
		Component: "gpu",
		Value:     "",
	}
	var sb strings.Builder

	sensor := fmt.Sprintf("/sys/class/hwmon/hwmon1/%s", c.Sensor)
	sensorValue, err := utils.GetSensorValue(sensor)
	if err != nil {
		log.Println("Error reading", sensor, err)
		return descriptor, err
	}

	tempCelcius:= utils.ReadSensorValue(sensorValue)
	sb.WriteString(fmt.Sprintf("GPU %5.1f °C", tempCelcius))
	descriptor.Value = sb.String()
	return descriptor, nil
}

func (c *GPU) Start(buffer chan descriptor.Descriptor) {
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

func (c *GPU) Stop() {
	c.Enabled.Store(false)
}

