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

const (
	Name  = "amdgpu"
	Label = "edge"
)

type GPU struct {
	Name       string
	Label      string
	SensorPath string
	Interval   time.Duration
	Enabled    atomic.Bool
}

// Attempt to create a sensor with data coming from the hwmon (name) and
// sensor (label). If either string is empty, internal default values will
// be assumed.
func New(name, label string, interval time.Duration) (*GPU, error) {
	gpu := &GPU{
		Interval: interval,
	}

	if len(name) == 0 {
		gpu.Name = Name
	} else {
		gpu.Name = name
	}

	if len(label) == 0 {
		gpu.Label = Label
	} else {
		gpu.Label = label
	}

	sensorPath, err := utils.FindSensorPath(gpu.Name, gpu.Label)
	if err != nil {
		return nil, err
	}

	gpu.SensorPath = sensorPath
	gpu.Update()
	return gpu, nil
}

func (g *GPU) Update() (descriptor.Descriptor, error) {
	log.Printf("INFO: Updating GPU temperature sensor=%s", g.Label)
	descriptor := descriptor.Descriptor{
		Component: "gpu",
		Value:     "",
	}
	var sb strings.Builder

	sensorValue, err := utils.GetSensorValue(g.SensorPath)
	if err != nil {
		log.Printf("ERROR: GetSensorValue sensor=%s err=%s", g.SensorPath, err)
		return descriptor, err
	}

	tempCelcius := utils.ReadSensorValue(sensorValue)
	sb.WriteString(fmt.Sprintf("GPU %5.1f °C", tempCelcius))
	descriptor.Value = sb.String()
	return descriptor, nil
}

func (g *GPU) Start(buffer chan descriptor.Descriptor) {
	g.Enabled.Store(true)

	go func() {
		for g.Enabled.Load() {
			descriptor, err := g.Update()
			if err != nil {
				log.Printf("ERROR: Update err=%s", err)
			} else {
				buffer <- descriptor
			}
			time.Sleep(g.Interval)
		}
	}()
}

func (g *GPU) Stop() {
	g.Enabled.Store(false)
}
