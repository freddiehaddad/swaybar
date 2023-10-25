package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

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
		log.Println("Error reading", sensor, err)
		return 0, err
	}
	sensorValueTrimmed := strings.TrimSuffix(string(sensorValueRaw), "\n")
	sensorValue, err := strconv.ParseInt(sensorValueTrimmed, 10, 32)
	if err != nil {
		log.Println("Error parsing value", sensorValueTrimmed, err)
		return 0, err
	}
	return sensorValue, nil
}

// Converts a value returned by GetSensorValue into a float64 temperature
// in degrees Celcius.
func ReadSensorValue(value int64) float64 {
	return float64(value) / 1000.0
}
