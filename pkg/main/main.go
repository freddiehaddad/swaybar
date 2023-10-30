package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	cputemp "github.com/freddiehaddad/swaybar/pkg/cpu/temp"
	cpuutil "github.com/freddiehaddad/swaybar/pkg/cpu/utilization"
	"github.com/freddiehaddad/swaybar/pkg/date"
	"github.com/freddiehaddad/swaybar/pkg/descriptor"
	"github.com/freddiehaddad/swaybar/pkg/gpu"
	"github.com/freddiehaddad/swaybar/pkg/interfaces"
	"github.com/freddiehaddad/swaybar/pkg/network"
	"gopkg.in/yaml.v3"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(os.Stderr)
}

type Component struct {
	Interval string `yaml:"interval"`
	Device   string `yaml:"device"`
	Sensor   string `yaml:"sensor"`
	Format   string `yaml:"format"`
	Order    int    `yaml:"order"`
}

func ParseInterval(interval string) (time.Duration, error) {
	if len(interval) == 0 {
		return time.Second, nil
	}

	parsed, err := time.ParseDuration(interval)
	if err != nil {
		log.Println("Error parsing interval", interval, err)
	}
	return parsed, err
}

func GenerateRenderOrder(components map[string]Component, order []string) error {
	for component, settings := range components {
		order[settings.Order] = component
	}
	return nil
}

func main() {
	// load the config file
	configFile := "./config/config.yml"
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalln("Error reading config file", configFile, err)
	}

	// parse yaml format
	componentConfigs := map[string]Component{}
	err = yaml.Unmarshal(configBytes, &componentConfigs)
	if err != nil {
		log.Fatalln("Error parsing config", configFile, err)
	}

	// order to render components on the status bar -- will come from a config file
	renderOrder := make([]string, len(componentConfigs))
	err = GenerateRenderOrder(componentConfigs, renderOrder)
	if err != nil {
		log.Fatalln("Error generating render order", err)
	}

	// list of components
	components := map[string]interfaces.Runnable{}

	// component updates arrive via a buffered channel asynchronously
	componentUpdates := make(chan descriptor.Descriptor, len(componentConfigs))

	// last received component update
	statusBar := map[string]descriptor.Descriptor{}

	// create the components
	for component, settings := range componentConfigs {
		switch component {
		case "cputemp":
			log.Println("Creating cpu temperature component")
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Println("Failed to parse interval", interval, err, "using default value of 1s")
				interval = time.Second
			}
			cputemp, _ := cputemp.New(settings.Sensor, interval)
			components["cputemp"] = cputemp
		case "cpuutil":
			log.Println("Creating cpu utilization component")
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Println("Failed to parse interval", interval, err, "using default value of 1s")
				interval = time.Second
			}
			cpuutil, _ := cpuutil.New(interval)
			components["cpuutil"] = cpuutil
		case "gpu":
			log.Println("Creating gpu component")
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Println("Failed to parse interval", interval, err, "using default value of 1s")
				interval = time.Second
			}
			gpu, _ := gpu.New(settings.Sensor, interval)
			components["gpu"] = gpu
		case "network":
			log.Println("Creating network component")
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Println("Failed to parse interval", interval, err, "using default value of 1s")
				interval = time.Second
			}
			network, _ := network.New(settings.Device, interval)
			components["network"] = network
		case "date":
			log.Println("Creating date component")
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Println("Failed to parse interval", interval, err, "using default value of 1s")
				interval = time.Second
			}
			date, _ := date.New(settings.Format, interval)
			components["date"] = date
		default:
			log.Println("Unknown component: ", component)
		}
	}

	// start the components
	for name, component := range components {
		log.Println("Starting", name)
		component.Start(componentUpdates)
	}

	// render the statusbar when updates arrive
	for descriptor := range componentUpdates {
		// store the update
		log.Println("Update from", descriptor.Component)
		statusBar[descriptor.Component] = descriptor

		// generate the statusbar
		sep := ""
		stringBuilder := strings.Builder{}
		for _, component := range renderOrder {
			descriptor, exists := statusBar[component]

			// unless we wait for all components to generate the
			// first update, we might have empty components
			if !exists {
				continue
			}

			stringBuilder.WriteString(sep)
			stringBuilder.WriteString(descriptor.Value)
			sep = " | "
		}

		fmt.Println(stringBuilder.String())
	}
}
