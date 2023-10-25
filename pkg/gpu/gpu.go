package gpu

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/freddiehaddad/swaybar/pkg/descriptor"
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

func getCurrentTemp(sensor string) (float64, error) {
	dataPath := fmt.Sprintf("/sys/class/hwmon/hwmon1/%s", sensor)
	dataRaw, err := os.ReadFile(dataPath)
	if err != nil {
		log.Println("Error reading", dataPath, err)
		return 0, err
	}
	dataString := strings.TrimSuffix(string(dataRaw), "\n")
	tempRaw, err := strconv.ParseInt(dataString, 10, 64)
	if err != nil {
		log.Println("Error parsing int64", dataString, err)
		return 0, err
	}
	tempCelcius := float64(tempRaw) / 1000.0
	return tempCelcius, nil

}

func (c *GPU) Update() (descriptor.Descriptor, error) {
	log.Println("Updating", c.Sensor)
	descriptor := descriptor.Descriptor{
		Component: "gpu",
		Value:     "",
	}
	var sb strings.Builder

	tempCelcius, err := getCurrentTemp(c.Sensor)
	if err != nil {
		log.Println("Error reading", c.Sensor, err)
		return descriptor, err
	}

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

