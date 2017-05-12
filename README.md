# ledctl - LED device reporting and setting tool for Linux, written in golang.

This is a rather plain golang-based method of gathering and setting LED data
from a wide range of devices, if-and-only-if supported by the Linux kernel.

When it comes to LED device data on Linux, I expect there could potentially
be some non-conformity (some devices may have odd behaviour), so go ahead
and shoot me an email if you see any obvious issues on your hardware.

My reason for writing this was a need for a more minimalist method of
accessing a number of customizable LED devices, many of which have their
own API which sometimes are heavy-weight implementations. As the Linux
kernel already controls these, a simple golang commandline interface
works for trivial situations.

Maybe one day it will be more fleshed out, but for now it is more of a
simple tool. Feel free to fork it and use it for other projects if you find
it useful.


# Requirements

The following is needed in order for this to function as intended:

* Linux kernel 4.4+
* golang

Older kernels could still give some kind of result, but I *think* most of
the newer Linux distros ought to contain all of the modules need to read
and write data to the brightness pseudo-files. Feel free to email me if
this is incorrect.


# Running

0) Build this program as you would a simple golang module.

    go build ledctl

1) Run the program to gather a list of devices.

    ledctl

2) Set the brightness level of a given device (this requires sudo).

    sudo ledctl --device input2::scrolllock --level 1


# Author

Written by Robert Bisewski at Ibis Cybernetics. For more information, contact:

* Website -> www.ibiscybernetics.com

* Email -> contact@ibiscybernetics.com
