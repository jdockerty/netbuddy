package main

import (
	"fmt"
	"testing"
)

func TestGetAddressCount(t *testing.T) {
	_, ipnet := parseIPInfo("192.168.1.1/24")
	addressCount := getAddressCount(ipnet)
	fmt.Printf("Total address count for %s is %d\n", ipnet, addressCount)
	if addressCount != 256 {
		t.Errorf("/24 CIDR returned %d, required 256.", addressCount)
	}
}

func TestGetCommonPorts(t *testing.T) {
	serviceTest := "SSH"
	response := getCommonPorts(serviceTest)

	if response.transportLayerProtocol != "TCP" {
		t.Errorf("Transport protocol of SSH is TCP. Returned %s.", response.transportLayerProtocol)
	}

	if response.extraInfoLink != "https://en.wikipedia.org/wiki/Secure_Shell" {
		t.Errorf("Invalid link returned from SSH service input. Returned %s.", response.extraInfoLink)
	}

	if response.commonPortNumbers[0] != 22 {
		t.Errorf("Incorrect port returned for SSH, wanted 22. Returned %d", response.commonPortNumbers[0])
	}
	fmt.Printf("Transport protocol for SSH is %s\n", response.transportLayerProtocol)
	fmt.Printf("Port number for SSH is %d\n", response.commonPortNumbers[0])
	fmt.Printf("Wiki link is %s\n", response.extraInfoLink)
}

func TestGetCommonPortsNonSupported(t *testing.T) {
	serviceTest := "a non existent service" // Incorrect name to test incorrect/unsupported services
	response := getCommonPorts(serviceTest)

	if response.err != 1 {
		t.Errorf("Error number was not set. Returned %+v\n", response)
	}

}

func TestGetSubnetInfo(t *testing.T) {
	_, ipnet := parseIPInfo("192.168.1.1/23")
	response := getSubnetData(ipnet)
	if response.networkAddress.String() != "192.168.0.0" {
		t.Errorf("Invalid network address. Expected 192.168.0.0, but returned %s\n", response.networkAddress)
	}
}

func TestShowIPv4Range(t *testing.T) {
	if shouldReturnTrue := ipv4PrivateAddressRange(); shouldReturnTrue != true {
		t.Errorf("Showing the IPv4 range command did not execute successfully.")
	}
}

func TestShowCmdHelp(t *testing.T) {
	if shouldReturnTrue := showCmdHelp(); shouldReturnTrue != true {
		t.Errorf("Displaying the help subcommand of 'show' was executed successfully.")
	}
}

func TestShowInterfaceData(t *testing.T) {
	if shouldReturnTrue := showinterfaceData(); shouldReturnTrue != true {
		t.Errorf("Failed to display network interface information.")
	}
}

func TestConvertToCIDR(t *testing.T) {
	testSubnetMask := "255.255.255.192"
	result := convertToCIDRNotation(testSubnetMask)

	if result != 26 {
		t.Errorf("Incorrect mask returned. Expected /26, but returned /%d", result)
	}
}