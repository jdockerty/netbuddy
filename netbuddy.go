package main

import (
	"fmt"
	"net"
	"github.com/apparentlymart/go-cidr/cidr"
	"strings"
)

type SubnetInfo struct {
	networkAddress net.IP
	broadcastAddress net.IP
	firstUsuableAddress net.IP
	lastUsuableAddress net.IP
	totalAddressCount uint64
}
func parseIPInfo(ipString string) (net.IP, *net.IPNet) {
	ip, ipnet, err := net.ParseCIDR(ipString)
	if err != nil {
		panic(err)
	}
	return ip, ipnet
}

func getAddressCount(ipnet *net.IPNet) uint64 {
	return cidr.AddressCount(ipnet)
}

func getNetworkAndBroadcast(ipnet *net.IPNet) (net.IP, net.IP) {
	networkAddress, broadcastAddress := cidr.AddressRange(ipnet)
	return networkAddress, broadcastAddress
}

func getSubnetInfo(ipnet *net.IPNet) SubnetInfo {
	networkAddr, broadcastAddr := getNetworkAndBroadcast(ipnet)
	firstUsuableAddr := cidr.Inc(networkAddr)
	lastUsuableAddr := cidr.Dec(broadcastAddr)
	totalAddresses := getAddressCount(ipnet)

	subnetInfoResponse := SubnetInfo{
		networkAddress: networkAddr,
		broadcastAddress: broadcastAddr,
		firstUsuableAddress: firstUsuableAddr,
		lastUsuableAddress: lastUsuableAddr,
		totalAddressCount: totalAddresses,
	}
	
	return subnetInfoResponse
}

func getCommonPorts(service string) (int, string) {
	var transportProtocol string

	// if dns is used, reassign to domain string and search using net.LookupPort inbuilt function.
	if strings.ToLower(service) == "dns" {
		service = "domain"

		portNum, err := net.LookupPort("udp", service)
		if err != nil {
			panic(err)
		}
		transportProtocol = "UDP"

		fmt.Println(transportProtocol, portNum)

		return portNum, transportProtocol
	} else {
		portNum, err := net.LookupPort("tcp", service)
		if err != nil {
			panic(err)
		}
		transportProtocol = "TCP"

		fmt.Println(transportProtocol, portNum)

		return portNum, transportProtocol
	}

}

func main() {
	myIP := "192.168.1.5/20"
	ip, ipnet := parseIPInfo(myIP)
	fmt.Println("IP",ip, ipnet)
	getSubnetInfo(ipnet)
	getCommonPorts("http")
	addressRangeReminder()
}