package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/freddiehaddad/swaybar/pkg/date"
	"github.com/freddiehaddad/swaybar/pkg/descriptor"
	"github.com/freddiehaddad/swaybar/pkg/interfaces"
	"github.com/freddiehaddad/swaybar/pkg/network"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(os.Stderr)
}

func main() {
	// list of components -- will come from a config file
	components := map[string]interfaces.Runnable{}

	// order to render components on the status bar -- will come from a config file
	renderOrder := []string{"network", "date"}

	// component updates arrive via a buffered channel asynchronously
	componentUpdates := make(chan descriptor.Descriptor, len(renderOrder))

	// last received component update
	statusBar := map[string]descriptor.Descriptor{}

	// create the components
	for _, component := range renderOrder {
		switch component {
		case "network":
			log.Println("Creating network component")
			network, _ := network.New("enp6s0", time.Second)
			components["network"] = network
		case "date":
			log.Println("Creating date component")
			date, _ := date.New(time.RFC850, time.Second)
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
