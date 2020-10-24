package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	Version      = "0.0"
	printVersion = false

	givenDevice     string
	brightnessLevel int

	ledsDirectory             = "/sys/class/leds/"
	brightnessFileLocation    = "brightness"
	maxBrightnessFileLocation = "max_brightness"
)

func init() {

	flag.BoolVar(&printVersion, "version", false,
		"Print the current version of this program and exit.")

	flag.StringVar(&givenDevice, "device", "",
		"LED Device; e.g. 'input2::scrolllock' ")

	flag.IntVar(&brightnessLevel, "level", -1,
		"Requested brightness level for LED device.")
}

func main() {

	output := ""
	var err error

	flag.Parse()

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
	} else {
		flag.Usage()
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf(output)

	os.Exit(0)
}
