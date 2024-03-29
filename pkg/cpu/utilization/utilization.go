package utilization

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

// See: man proc for /proc/stat column value meanings
const (
	user = iota
	nice
	system
	idle
	iowait
	irq
	softirq
	steal
	guest
	guest_nice
	num_fields
)

type CPUUtilization struct {
	Interval        time.Duration
	PrevStatsValues []int64
	Enabled         atomic.Bool
}

func New(interval time.Duration) (*CPUUtilization, error) {
	cpu := &CPUUtilization{
		Interval:        interval,
		PrevStatsValues: make([]int64, num_fields),
	}

	cpu.Update()
	return cpu, nil
}

func (c *CPUUtilization) Update() (descriptor.Descriptor, error) {
	log.Printf("INFO: Updating CPU utilization")
	descriptor := descriptor.Descriptor{
		Component: "cpuutil",
		Value:     "",
	}
	var sb strings.Builder

	statPath := "/proc/stat"
	statRaw, err := os.ReadFile(statPath)
	if err != nil {
		log.Printf("ERROR: ReadFile statPath=%s err=%s", statPath, err)
		return descriptor, err
	}

	currentStatValuesRaw, err := getStatValues(statRaw)
	if err != nil {
		log.Println(err)
		return descriptor, err
	}

	currentStatValues, err := parseInts(currentStatValuesRaw)
	if err != nil {
		log.Printf("ERROR: parseInts currentStatValuesRaw=%v err=%s", currentStatValuesRaw, err)
		return descriptor, err
	}

	previousStatValuesSum := sumArray(c.PrevStatsValues)
	currentStatValuesSum := sumArray(currentStatValues)

	previousIdleValue := c.PrevStatsValues[idle]
	currentIdleValue := currentStatValues[idle]

	utilizationDelta := currentStatValuesSum - previousStatValuesSum
	idleDelta := currentIdleValue - previousIdleValue

	netUtilizationDelta := utilizationDelta - idleDelta

	cpuUtilization := 100.0 * float64(netUtilizationDelta) / float64(utilizationDelta)

	c.PrevStatsValues = currentStatValues

	sb.WriteString(fmt.Sprintf("CPU %5.1f%%", cpuUtilization))
	descriptor.Value = sb.String()
	return descriptor, nil
}

func (c *CPUUtilization) Start(buffer chan descriptor.Descriptor) {
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

func (c *CPUUtilization) Stop() {
	c.Enabled.Store(false)
}

func getStatValues(bytes []byte) ([]string, error) {
	const expectedLength = 2
	const numSplits = 2
	const rawSeparator = "\n"
	const valueSeparator = " "
	const valuesExpected = 10
	const firstValue = "cpu"

	finalValues := []string{}

	s := string(bytes)
	split := strings.SplitAfterN(s, rawSeparator, numSplits)
	if len(split) != 2 {
		err := fmt.Errorf("length error len=%d got=%d s=%s rawSeparator=%s numSplits=%d", expectedLength, len(split), s, rawSeparator, numSplits)
		return finalValues, err
	}

	stats := split[0]

	if len(stats) <= len(firstValue) {
		err := fmt.Errorf("length error stats=%s len=%d <= firstValue=%s len=%d", stats, len(stats), firstValue, len(firstValue))
		return finalValues, err
	}

	stats = strings.TrimPrefix(stats, firstValue)
	stats = strings.TrimSpace(stats)

	values := strings.Split(stats, valueSeparator)
	if len(values) != valuesExpected {
		err := fmt.Errorf("split error stats=%s did not produce valuesExpected=%d got=%d", stats, valuesExpected, len(values))
		return finalValues, err
	}

	return values, nil
}

func sumArray(values []int64) int64 {
	sum := int64(0)

	for _, value := range values {
		sum += value
	}

	return sum
}

func parseInts(values []string) ([]int64, error) {
	intValues := make([]int64, len(values))
	for index, value := range values {
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return intValues, err
		}
		intValues[index] = intValue
	}
	return intValues, nil
}
