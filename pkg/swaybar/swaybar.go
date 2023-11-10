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
	Order    int    `yaml:"order"`
	Device   string `yaml:"device"`
	Name     string `yaml:"name"`
	Label    string `yaml:"label"`
	Format   string `yaml:"format"`
}

func ParseInterval(interval string) (time.Duration, error) {
	if len(interval) == 0 {
		return time.Second, nil
	}

	parsed, err := time.ParseDuration(interval)
	if err != nil {
		log.Printf("ERROR: ParseInterval interval=%s err=%s\n", interval, err)
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
		log.Fatalf("ERROR: ReadFile configFile=%s err=%s\n", configFile, err)
	}

	// parse yaml format
	componentConfigs := map[string]Component{}
	err = yaml.Unmarshal(configBytes, &componentConfigs)
	if err != nil {
		log.Fatalf("ERROR: Unmarshal configFile=%s err=%s\n", configFile, err)
	}

	// order to render components on the status bar -- will come from a config file
	renderOrder := make([]string, len(componentConfigs))
	err = GenerateRenderOrder(componentConfigs, renderOrder)
	if err != nil {
		log.Fatalf("ERROR: GenerateRenderOrder err=%s\n", err)
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
			log.Printf("INFO: Creating component=%q settings=%+v\n", component, settings)
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Fatalf("ERROR: ParseInterval interval=%s err=%s\n", settings.Interval, err)
			}
			cputemp, err := cputemp.New(settings.Name, settings.Label, interval)
			if err != nil {
				log.Fatalf("ERROR: Failed to create component=%q err=%q\n", component, err)
			}
			components["cputemp"] = cputemp
		case "cpuutil":
			log.Printf("INFO: Creating component=%q settings=%+v\n", component, settings)
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Fatalf("ERROR: ParseInterval interval=%s err=%s\n", settings.Interval, err)
			}
			cpuutil, _ := cpuutil.New(interval)
			components["cpuutil"] = cpuutil
		case "gpu":
			log.Printf("INFO: Creating component=%q settings=%+v\n", component, settings)
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Fatalf("ERROR: ParseInterval interval=%s err=%s\n", settings.Interval, err)
			}
			gpu, err := gpu.New(settings.Name, settings.Label, interval)
			if err != nil {
				log.Fatalf("ERROR: Failed to create component=%q err=%q\n", component, err)
			}
			components["gpu"] = gpu
		case "network":
			log.Printf("INFO: Creating component=%q settings=%+v\n", component, settings)
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Fatalf("ERROR: ParseInterval interval=%s err=%s\n", settings.Interval, err)
			}
			network, _ := network.New(settings.Device, interval)
			components["network"] = network
		case "date":
			log.Printf("INFO: Creating component=%q settings=%+v\n", component, settings)
			interval, err := ParseInterval(settings.Interval)
			if err != nil {
				log.Fatalf("ERROR: ParseInterval interval=%s err=%s\n", settings.Interval, err)
			}
			date, _ := date.New(settings.Format, interval)
			components["date"] = date
		default:
			log.Fatalf("ERROR: Unknown component=%s settings=%+v\n", component, settings)
		}
	}

	// start the components
	for name, component := range components {
		log.Printf("INFO: Starting name=%s\n", name)
		component.Start(componentUpdates)
	}

	// render the statusbar when updates arrive
	for descriptor := range componentUpdates {
		// store the update
		log.Printf("INFO: Update from descriptor=%+v\n", descriptor)
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
