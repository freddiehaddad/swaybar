package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const SensorPath = "/sys/class/hwmon"

// Returns the sensor value in its raw form as exposed by the /sys file system.
// For example, with the Intel Core i9-12900k, the temperature will be returned
// as an integer value in degrees Celcius. Fractional parts of the temperature
// (if supported by the CPU) will require conversion.  Simply put, if the
// temperature is 20.5 degrees Celcius, the value returned by this function will
// be 20500.
//
// The sensor string must be a valid full path to the sensor:
//
//   - /sys/class/hwmon/hwmon1/temp1_input
//
// Any error that prevents reading the value will be returned.
func GetSensorValue(sensor string) (int64, error) {
	sensorValueRaw, err := os.ReadFile(sensor)
	if err != nil {
		return 0, err
	}
	sensorValueTrimmed := strings.TrimSuffix(string(sensorValueRaw), "\n")
	sensorValue, err := strconv.ParseInt(sensorValueTrimmed, 10, 32)
	if err != nil {
		return 0, err
	}
	return sensorValue, nil
}

// Converts a value returned by GetSensorValue into a float64 temperature
// in degrees Celcius.
func ReadSensorValue(value int64) float64 {
	return float64(value) / 1000.0
}

// Returns true if the mode has the symbolic link flag set.
func symbolicLink(mode fs.FileMode) bool {
	return mode&fs.ModeSymlink != 0
}

// Searches all directories under SensorPath, following symbolic links, for
// an hwmonX/name file whose contents matches name and a *_input file whose
// corresponding *_label file matches label.  Returns the full path to the
// matching *_input file containing the sensor data, or an error if no match
// was was found.
func FindSensorPath(name, label string) (string, error) {
	sensorPath := ""
	entries, err := os.ReadDir(SensorPath)
	if err != nil {
		return sensorPath, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return sensorPath, err
		}

		if !symbolicLink(info.Mode()) {
			err := errors.New("Expected " + entry.Name() + " to be a symbolic link.")
			return sensorPath, err
		}

		fp := filepath.Join(SensorPath, entry.Name())
		link, err := os.Readlink(fp)
		if err != nil {
			return sensorPath, err
		}

		path := filepath.Join(SensorPath, link)
		entries, err := os.ReadDir(path)
		if err != nil {
			return sensorPath, err
		}

		nameFound := false
		labelFound := false
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), "name") {
				fp := filepath.Join(path, entry.Name())
				raw, err := os.ReadFile(fp)
				if err != nil {
					return sensorPath, err
				}

				data := string(raw)
				data = strings.TrimSpace(data)

				if data == name {
					nameFound = true
				}
			} else if strings.HasSuffix(entry.Name(), "label") { // e.g. temp1_label
				fp := filepath.Join(path, entry.Name())
				raw, err := os.ReadFile(fp)
				if err != nil {
					return sensorPath, err
				}

				data := string(raw)
				data = strings.TrimSpace(data)
				data = strings.ToLower(data)

				if !strings.HasPrefix(data, label) {
					continue
				}

				prefix := strings.TrimSuffix(entry.Name(), "label")
				filename := fmt.Sprintf("%s%s", prefix, "input")
				fp = filepath.Join(path, filename)

				_, err = os.ReadFile(fp)
				if err != nil {
					return sensorPath, err
				}

				labelFound = true
				sensorPath = fp
			}

			if nameFound && labelFound {
				return sensorPath, nil
			}
		}
	}

	err = errors.New("No sensor data found for " + name + "/" + label)
	return sensorPath, err
}
