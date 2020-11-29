# NetBuddy
<img src="https://travis-ci.com/jdockerty/netbuddy.svg?token=xPjFq5JeCTp415MsJdAD&branch=master">

Simple CLI tool for networking related information:
* Addresses contained within a network
* Total address count for a particular subnet mask.
* Showing common ports for services
* etc...

Examples are shown below.

## Installation

### With Go installed
Executing the commands below will download the repo onto your machine, build the binary using Go, and then move the binary into your `bin` folder for direct execution from the command-line.

```
git clone https://github.com/jdockerty/netbuddy.git
go build
sudo mv netbuddy /usr/local/bin
```

## Examples

Running the command `netbuddy subnet -display <IP/X>` will show the network, broadcast, first and last assignable addresses for the subnet. <br><br>
<img src="https://github.com/jdockerty/netbuddy/blob/master/READMEimages/displayExample.png">

Another example is also for displaying the common services or protocol, such as BGP. <br><br>
<img src="https://github.com/jdockerty/netbuddy/blob/master/READMEimages/showBGPExample.png">


