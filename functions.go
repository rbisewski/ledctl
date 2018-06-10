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

	// Holds the string value to return.
	output := ""

	// Print out a few lines telling the user that the program has started.
	output += "\nDevice Name\t\t\tBrightness\tMaximum Brightness\n\n"

	// Attempt to read in our file contents.
	listOfDeviceDirs, err := ioutil.ReadDir(ledsDirectory)

	// If an error occurs, end the program.
	if err != nil {
		return "", err
	}

	// For each of the devices...
	for _, dir := range listOfDeviceDirs {

		// Determine the brightness and max brightness files.
		brightnessFile := ledsDirectory + dir.Name() + "/" +
			brightnessFileLocation
		maxBrightnessFile := ledsDirectory + dir.Name() + "/" +
			maxBrightnessFileLocation

		// Attempt to obtain the brightness
		brightnessAsBytes, err := ioutil.ReadFile(brightnessFile)

		// If an error occurs, go back.
		if err != nil {
			return "", err
		}

		// Attempt to convert the brightness value as bytes into an integer.
		brightness, err := strconv.Atoi(
			strings.Trim(string(brightnessAsBytes), " \n"))

		// If an error occurs, skip to the next device.
		if err != nil {
			return "", err
		}

		// Attempt to obtain the brightness
		maxBrightnessAsBytes, err := ioutil.ReadFile(maxBrightnessFile)

		// If an error occurs, skip to the next device.
		if err != nil {
			return "", err
		}

		// Attempt to convert the brightness value as bytes into an integer.
		maxBrightness, err := strconv.Atoi(
			strings.Trim(string(maxBrightnessAsBytes), " \n"))

		// If an error occurs, skip to the next device.
		if err != nil {
			return "", err
		}

		// Trim away any whitespace from the device name.
		name := strings.Trim(dir.Name(), " \n\t\v")

		// Finally print out the current line.
		output += name + " \t\t" + strconv.Itoa(brightness) + "\t\t" +
			strconv.Itoa(maxBrightness) + "\n"
	}

	// Append a newline to the end of the output.
	output += "\n"

	// If got to this point, no error has occurred.
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

	// input validation
	if len(device) < 1 || level < 0 {
		return "", fmt.Errorf("setLedBrightness() --> invalid input")
	}

	// Trim away any whitespace from the device name.
	device = strings.Trim(device, " \n\t\v")

	// Determine the brightness and max brightness files.
	brightnessFile := ledsDirectory + device + "/" +
		brightnessFileLocation
	maxBrightnessFile := ledsDirectory + device + "/" +
		maxBrightnessFileLocation

	// Attempt to obtain the max brightness as bytes.
	maxBrightnessAsBytes, err := ioutil.ReadFile(maxBrightnessFile)

	// If an error occurs, skip to the next device.
	if err != nil {
		return "", err
	}

	// Attempt to convert the brightness value as bytes into an integer.
	maxBrightness, err := strconv.Atoi(
		strings.Trim(string(maxBrightnessAsBytes), " \n"))

	// If an error occurs, skip to the next device.
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

	// Attempt to cast the new brightness level to a string.
	newBrightnessAsString := strconv.Itoa(level)

	// Attempt to write the new brightness level to the file.
	newBrightnessAsBytes := []byte(newBrightnessAsString)
	err = ioutil.WriteFile(brightnessFile, newBrightnessAsBytes, 0644)

	// If an error occurred, send it back.
	if err != nil {
		return "", err
	}

	// Assemble the success message.
	output := "The device [" + device + "] is now set to a brightness " +
		"level of [" + strconv.Itoa(level) + "]"

	// Send it back.
	return output, nil
}
