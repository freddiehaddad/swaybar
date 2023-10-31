package date

import (
	"sync/atomic"
	"log"
	"time"

	"github.com/freddiehaddad/swaybar/pkg/descriptor"
)

type Date struct {
	Format   string
	Interval time.Duration
	Enabled  atomic.Bool
}

func New(format string, interval time.Duration) (*Date, error) {
	if len(format) == 0 {
		format = time.RFC850

	}
	date := &Date{
		Format: format,
		Interval: interval,
	}

	date.Update()
	return date, nil
}

func (d *Date) Update() (descriptor.Descriptor, error) {
	log.Printf("INFO: Updating date")
	descriptor := descriptor.Descriptor{
		Component: "date",
		Value:     "",
	}

	date := time.Now().Format(d.Format)
	descriptor.Value = date

	return descriptor, nil
}

func (d *Date) Start(buffer chan descriptor.Descriptor) {
	go func() {
		d.Enabled.Store(true)

		for d.Enabled.Load() {
			descriptor, err := d.Update()
			if err != nil {
				log.Printf("ERROR: Update err=%s", err)
			} else {
				buffer <- descriptor
			}
			time.Sleep(d.Interval)
		}
	}()
}

func (d *Date) Stop() {
	d.Enabled.Store(false)
}
