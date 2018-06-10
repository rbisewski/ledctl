package main

//
// Imports
//
import (
	"flag"
	"fmt"
	"os"
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
