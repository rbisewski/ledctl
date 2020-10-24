package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//! Function that obtains LED info.
/*
 * @return    string     output
 * @return    error      error message, if any
 */
func getLedInfo() (string, error) {

	output := "\nDevice Name\t\t\tBrightness\tMaximum Brightness\n\n"

	listOfDeviceDirs, err := ioutil.ReadDir(ledsDirectory)
	if err != nil {
		return "", err
	}

	for _, dir := range listOfDeviceDirs {

		// Determine the brightness and max brightness files.
		brightnessFile := ledsDirectory + dir.Name() + "/" +
			brightnessFileLocation
		maxBrightnessFile := ledsDirectory + dir.Name() + "/" +
			maxBrightnessFileLocation

		brightnessAsBytes, err := ioutil.ReadFile(brightnessFile)
		if err != nil {
			return "", err
		}

		brightness, err := strconv.Atoi(strings.Trim(string(brightnessAsBytes), " \n"))
		if err != nil {
			return "", err
		}

		maxBrightnessAsBytes, err := ioutil.ReadFile(maxBrightnessFile)
		if err != nil {
			return "", err
		}

		maxBrightness, err := strconv.Atoi(strings.Trim(string(maxBrightnessAsBytes), " \n"))
		if err != nil {
			return "", err
		}

		output += strings.Trim(dir.Name(), " \n\t\v") + " \t\t" + strconv.Itoa(brightness) + "\t\t" + strconv.Itoa(maxBrightness) + "\n"
	}

	output += "\n"

	return output, nil
}

//! Function to set the brightness level of a given LED.
/*
 * @param     string    given LED device
 * @param     int       requested brightness level
 *
 * @return    string    response message, if any
 * @return    error     error message, if any
 */
func setLedBrightness(device string, level int) (string, error) {

	if len(device) < 1 || level < 0 {
		return "", fmt.Errorf("setLedBrightness() --> invalid input")
	}

	device = strings.Trim(device, " \n\t\v")

	// Determine the brightness and max brightness files.
	brightnessFile := ledsDirectory + device + "/" + brightnessFileLocation
	maxBrightnessFile := ledsDirectory + device + "/" + maxBrightnessFileLocation

	maxBrightnessAsBytes, err := ioutil.ReadFile(maxBrightnessFile)
	if err != nil {
		return "", err
	}

	maxBrightness, err := strconv.Atoi(strings.Trim(string(maxBrightnessAsBytes), " \n"))
	if err != nil {
		return "", err
	}

	// Ensure that the requested brightness level does not exceed the
	// maximum brightness level possible by the device.
	if level > maxBrightness {
		return "", fmt.Errorf("Error: Requested brightness of (%d) is "+
			"beyond the maximum possible of the device (%d).", level,
			maxBrightness)
	}

	newBrightnessAsString := strconv.Itoa(level)
	newBrightnessAsBytes := []byte(newBrightnessAsString)

	err = ioutil.WriteFile(brightnessFile, newBrightnessAsBytes, 0644)
	if err != nil {
		return "", err
	}

	return "The device [" + device + "] is now set to a brightness level of [" + strconv.Itoa(level) + "]", nil
}
