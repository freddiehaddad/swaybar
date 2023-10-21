package interfaces

import "github.com/freddiehaddad/swaybar/pkg/descriptor"

type Runnable interface {
	Start(buffer chan descriptor.Descriptor)
	Stop()
}
