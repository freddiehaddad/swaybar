package network

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

const (
	gbps = 1 << 30
	mbps = 1 << 20
	kbps = 1 << 10
)

type Network struct {
	Device        string
	PrevRxBytes   int64
	PrevTxBytes   int64
	PrevTimestamp int64
	Interval      time.Duration
	Enabled       atomic.Bool
}

func New(device string, interval time.Duration) (*Network, error) {
	network := &Network{
		Device:        device,
		PrevRxBytes:   0,
		PrevTxBytes:   0,
		PrevTimestamp: 0,
		Interval:      interval,
	}

	network.Update()
	return network, nil
}

func getBytesTransferred(device, direction string) (int64, error) {
	dataPath := fmt.Sprintf("/sys/class/net/%s/statistics/%s", device, direction)
	dataRaw, err := os.ReadFile(dataPath)
	if err != nil {
		log.Println("Error reading", dataPath, err)
		return 0, err
	}
	dataString := strings.TrimSuffix(string(dataRaw), "\n")
	bytesTransferred, err := strconv.ParseInt(dataString, 10, 64)
	if err != nil {
		log.Println("Error parsing int64", dataString, err)
		return 0, err
	}
	return bytesTransferred, nil
}

func convertBytesToBits(bytes int64) int64 {
	return bytes * 8
}

func convertNanosecondsToSeconds(nanoseconds int64) float64 {
	return float64(nanoseconds) / float64(time.Second)
}

func shortenThroughput(bitsPerSecond float64) (float64, string) {
	if bitsPerSecond >= gbps {
		gbpsThroughput := bitsPerSecond / float64(gbps)
		return gbpsThroughput, "Gbps"
	} else if bitsPerSecond >= mbps {
		mbpsThroughput := bitsPerSecond / float64(mbps)
		return mbpsThroughput, "Mbps"
	} else if bitsPerSecond >= kbps {
		kbpsThroughput := bitsPerSecond / float64(kbps)
		return kbpsThroughput, "Kbps"
	} else {
		return bitsPerSecond, "bps"
	}
}

func calculateThroughput(prevTimeNanoseconds, prevBytes, currTimeNanoseconds, currBytes int64) (float64, string) {
	timeElapsedNanoseconds := currTimeNanoseconds - prevTimeNanoseconds
	timeElapsedSeconds := convertNanosecondsToSeconds(timeElapsedNanoseconds)

	bytesTransferred := currBytes - prevBytes
	bitsTransferred := convertBytesToBits(bytesTransferred)

	throughput := float64(bitsTransferred) / timeElapsedSeconds
	rate, unit := shortenThroughput(throughput)
	return rate, unit
}

func (n *Network) Update() (descriptor.Descriptor, error) {
	log.Println("Updating", n.Device)
	var s string
	descriptor := descriptor.Descriptor{
		Component: "network",
		Value:     "",
	}
	var sb strings.Builder

	currTimestamp := time.Now().UnixNano()

	s = "rx_bytes"
	rxBytesTransferred, rxErr := getBytesTransferred(n.Device, s)
	if rxErr != nil {
		log.Println("Error getting", s, rxErr)
		return descriptor, rxErr

	}

	s = "tx_bytes"
	txBytesTransferred, txErr := getBytesTransferred(n.Device, s)
	if txErr != nil {
		log.Println("Error getting", s, txErr)
		return descriptor, txErr
	}

	rxThroughput, rxUnit := calculateThroughput(n.PrevTimestamp, n.PrevRxBytes, currTimestamp, rxBytesTransferred)
	txThroughput, txUnit := calculateThroughput(n.PrevTimestamp, n.PrevTxBytes, currTimestamp, txBytesTransferred)

	n.PrevRxBytes = rxBytesTransferred
	n.PrevTxBytes = txBytesTransferred
	n.PrevTimestamp = currTimestamp

	sb.WriteString(fmt.Sprintf("%2s", "D"))
	sb.WriteString(fmt.Sprintf("%8.02f", rxThroughput))
	sb.WriteString(fmt.Sprintf("%5s", rxUnit))

	sb.WriteString(fmt.Sprintf("%2s", "U"))
	sb.WriteString(fmt.Sprintf("%8.02f", txThroughput))
	sb.WriteString(fmt.Sprintf("%5s", txUnit))

	descriptor.Value = sb.String()
	return descriptor, nil
}

func (n *Network) Start(buffer chan descriptor.Descriptor) {
	n.Enabled.Store(true)

	go func() {
		for n.Enabled.Load() {
			descriptor, err := n.Update()
			if err != nil {
				log.Println("Error during update", err)
			} else {
				buffer <- descriptor
			}
			time.Sleep(n.Interval)
		}
	}()
}

func (n *Network) Stop() {
	n.Enabled.Store(false)
}
