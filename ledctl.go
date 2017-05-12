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
    "fmt"
    "flag"
    "io/ioutil"
    "os"
    "strconv"
    "strings"
)

//
// Globals
//
var (

    // Requested LED device, if any
    givenDevice string

    // Requested brightness level for the given LED device, if any
    brightnessLevel int

    // Current location of the LED sensor data, as of kernel 4.4+
    leds_directory = "/sys/class/leds/"

    // Attribute file for storing the hardware device current brightness.
    brightness_file_location = "brightness"

    // Attribute file for storing the hardware device current maxbrightness.
    max_brightness_file_location = "max_brightness"
)

// Initialize the argument input flags.
func init() {

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
    var output string = ""
    var err error     = nil

    // Parse the flags, if any.
    flag.Parse()

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
    list_of_device_dirs, err := ioutil.ReadDir(leds_directory)

    // If an error occurs, end the program.
    if err != nil {
        return "", err
    }

    // For each of the devices...
    for _, dir := range list_of_device_dirs {

        // Determine the brightness and max brightness files.
        brightness_file     := leds_directory + dir.Name() + "/" +
          brightness_file_location
        max_brightness_file := leds_directory + dir.Name() + "/" +
          max_brightness_file_location

        // Attempt to obtain the brightness
        brightness_as_bytes, err := ioutil.ReadFile(brightness_file)

        // If an error occurs, go back.
        if err != nil {
            return "", err
        }

        // Attempt to convert the brightness value as bytes into an integer.
        brightness, err := strconv.Atoi(
          strings.Trim(string(brightness_as_bytes), " \n"))

        // If an error occurs, skip to the next device.
        if err != nil {
            return "", err
        }

        // Attempt to obtain the brightness
        max_brightness_as_bytes, err := ioutil.ReadFile(max_brightness_file)

        // If an error occurs, skip to the next device.
        if err != nil {
            return "", err
        }

        // Attempt to convert the brightness value as bytes into an integer.
        max_brightness, err := strconv.Atoi(
          strings.Trim(string(max_brightness_as_bytes), " \n"))

        // If an error occurs, skip to the next device.
        if err != nil {
            return "", err
        }

        // Trim away any whitespace from the device name.
        name := strings.Trim(dir.Name(), " \n\t\v")

        // Finally print out the current line.
        output += name + ">\t" + strconv.Itoa(brightness) + "\t" +
          strconv.Itoa(max_brightness) + "\n"
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
    brightness_file     := leds_directory + device + "/" +
      brightness_file_location
    max_brightness_file := leds_directory + device + "/" +
      max_brightness_file_location

    // Attempt to obtain the max brightness as bytes.
    max_brightness_as_bytes, err := ioutil.ReadFile(max_brightness_file)

    // If an error occurs, skip to the next device.
    if err != nil {
        return "", err
    }

    // Attempt to convert the brightness value as bytes into an integer.
    max_brightness, err := strconv.Atoi(
      strings.Trim(string(max_brightness_as_bytes), " \n"))

    // If an error occurs, skip to the next device.
    if err != nil {
        return "", err
    }

    // Ensure that the requested brightness level does not exceed the
    // maximum brightness level possible by the device.
    if level > max_brightness {
        return "", fmt.Errorf("Error: Requested brightness of (%d) is " +
          "beyond the maximum possible of the device (%d).", level,
          max_brightness)
    }

    // Attempt to cast the new brightness level to a string.
    new_brightness_as_string := strconv.Itoa(level)

    // Attempt to write the new brightness level to the file.
    new_brightness_as_bytes  := []byte(new_brightness_as_string)
    err = ioutil.WriteFile(brightness_file, new_brightness_as_bytes, 0644)

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
