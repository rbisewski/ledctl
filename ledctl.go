/*
 * LED Control Tool
 *
 * Description: A simple tool written in golang, for the purposes of
 *              checking and setting the brightness of system LEDs.
 *
 *              Specifically, this will work on kernel version 4.4+ or
 *              newer, depending on whether or not your given LED has
 *              a driver / module for its device.
 *
 * Author: Robert Bisewski <contact@ibiscybernetics.com>
 */

//
// Package
//
package main

//
// Imports
//
import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//
// Globals
//
var (
	// default version value
	Version = "0.0"

	// Whether or not to print the current version of the program
	printVersion = false

	// Requested LED device, if any
	givenDevice string

	// Requested brightness level for the given LED device, if any
	brightnessLevel int

	// Current location of the LED sensor data, as of kernel 4.4+
	ledsDirectory = "/sys/class/leds/"

	// Attribute file for storing the hardware device current brightness.
	brightnessFileLocation = "brightness"

	// Attribute file for storing the hardware device max brightness.
	maxBrightnessFileLocation = "max_brightness"
)

// Initialize the argument input flags.
func init() {

	// Version mode flag
	flag.BoolVar(&printVersion, "version", false,
		"Print the current version of this program and exit.")

	// LED device flag
	flag.StringVar(&givenDevice, "device", "",
		"LED Device; e.g. 'input2::scrolllock' ")

	// Brightness level flag
	flag.IntVar(&brightnessLevel, "level", -1,
		"Requested brightness level for LED device.")
}

//
// PROGRAM MAIN
//
func main() {

	// String variable to hold eventual output, as well error variable.
	output := ""
	var err error

	// Parse the flags, if any.
	flag.Parse()

	// if requested, go ahead and print the version; afterwards exit the
	// program, since this is all done
	if printVersion {
		fmt.Println("ledctl v" + Version)
		os.Exit(0)
	}

	// If a given device is mentioned and a sane brightness level, go
	// ahead and attempt to set the device to that brightness level.
	if len(givenDevice) > 0 && brightnessLevel >= 0 {
		output, err = setLedBrightness(givenDevice, brightnessLevel)

		// No device or brightness level mentioned? Then print all of the LED
		// device info on the system.
	} else if len(givenDevice) < 1 && brightnessLevel < 0 {
		output, err = getLedInfo()

		// If a device or brightness is mentioned, but not both, then print the
		// usage info, since the user has supplied bad input.
	} else {
		flag.Usage()
	}

	// If an error occurs, print it out.
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Otherwise print the output to stdout.
	fmt.Printf(output)

	// If all is well, we can return quietly here.
	os.Exit(0)
}

//! Function that obtains LED info.
/*
 * @return    string     output
 * @return    error      error message, if any
 */
func getLedInfo() (string, error) {

	// Holds the string value to return.
	output := ""

	// Print out a few lines telling the user that the program has started.
	output += "\n----------------------------------\n" +
		"LED Brightness Info Tool for Linux\n\n" +
		"The following info is displayed:\n\n" +
		"* Device Name\n" +
		"* Brightness\n" +
		"* Maximum Brightness\n" +
		"----------------------------------\n\n"

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
		output += name + ">\t" + strconv.Itoa(brightness) + "\t" +
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
