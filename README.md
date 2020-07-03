# NetBuddy
Simple CLI tool for networking related information:
* Addresses contained within a network
* Total address count for a particular subnet mask.
* Showing common ports for services
* etc...

Examples are shown below.

This project was also conducted for utilising a CI server to perform tests on new commits on the remote master branch, the choice was made for Travis CI.

## Installation

Executing the commands below will download the repo onto your machine, you can then move the pre-built binary into your bin folder for execution from the terminal.

```
git clone https://github.com/jdockerty/netbuddy.git
cd netbuddy/
sudo mv netbuddy /usr/local/bin
```

## Examples

Running the command `netbuddy subnet -display <IP/X>` will show the network, broadcast, first and last assignable addresses for the subnet. <br><br>
<img src="https://github.com/jdockerty/netbuddy/blob/master/READMEimages/displayExample.png">

Another example is also for displaying the common services or protocol, such as BGP. <br><br>
<img src="https://github.com/jdockerty/netbuddy/blob/master/READMEimages/showBGPExample.png">


## Go Testing and Travis CI
_This section is for my personal notes_

A small number of tests were completed in the `netbuddy_test.go`, this provides a standardised way to test the return variables from varying functions, ensuring they remain the same across adding other features or refactoring code.

Upon conducting the tests, this has helped in ironing out issues which were overlooked, such as testing the `show service` output with a capitalised input, the initial switch statement was evaluating a `strings.ToLower(var)` input, but the input itself had not been altered to always be lowercase to conform to the keys within the map. This was resolved after noticing it through the unit test response.

Travis CI was setup to provide a way in which automated tests can be continually conducted upon each new commit on the remote master branch. The tests which are run are those which have been written in the corresponding test file and the golangci-lint tool is executed to test for stylistic errors, bugs, and to enforce error checking.

The image below shows the most recent build outcome from running the tests provided.

<img src="https://travis-ci.com/jdockerty/netbuddy.svg?token=xPjFq5JeCTp415MsJdAD&branch=master">
