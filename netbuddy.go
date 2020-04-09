package main

import (
	"fmt"
	"net"
	"github.com/apparentlymart/go-cidr/cidr"
	"strings"
	"log"
)

type SubnetInfo struct {
	networkAddress net.IP
	broadcastAddress net.IP
	firstUsuableAddress net.IP
	lastUsuableAddress net.IP
	totalAddressCount uint64
}

type PortsInfo struct {
	commonPortNumbers []int
	transportLayerProtocol string
	extraInfoLink string
}
func parseIPInfo(ipString string) (net.IP, *net.IPNet) {
	ip, ipnet, err := net.ParseCIDR(ipString)
	if err != nil {
		log.Fatal(err)
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

func populatePortsInfo(portNums []int, transportProto, link string) PortsInfo {
	portsInfoResponse := PortsInfo{
		commonPortNumbers : portNums,
		transportLayerProtocol : transportProto,
		extraInfoLink : link,
	}
	return portsInfoResponse
}

func getWikiName(service string) string {

	abbreviationsToWikiName := map[string]string{
		"dns" : "Domain_Name_System",
		"dhcp" : "Dynamic_Host_Configuration_Protocol",
		"rdp" : "Remote_Desktop_Protocol",
		"smtp" : "Simple_Network_Management_Protocol",
		"ssh" : "Secure_Shell",
		"telnet" : "telnet",
		"ftp" : "File_Transfer_Protocol",
		"http" : "Hypertext_Transfer_Protocol",
		"https" : "HTTPS",
		"imap" : "Internet_Message_Access_Protocol",
		"pop3" : "Post_Office_Protocol",
		"ldap" : "Lightweight_Directory_Access_Protocol",
		"bgp" : "Border_Gateway_Protocol", 
	}

	return abbreviationsToWikiName[service]
}

func printPortInfo(portInfo PortsInfo) {
	fmt.Printf("Port Numbers: %d\nTransport Protocol(s): %s\nFor more information on this protocol visit %s\n",
	portInfo.commonPortNumbers, portInfo.transportLayerProtocol, portInfo.extraInfoLink)
}

func getCommonPorts(service string) PortsInfo {

	wikiString := "https://en.wikipedia.org/wiki/"

	switch strings.ToLower(service) {
	case "dns":
		info := populatePortsInfo([]int{53}, "UDP", wikiString + getWikiName(service))
		return info

	case "dhcp":
		info := populatePortsInfo([]int{67,68}, "UDP", wikiString + getWikiName(service))
		return info

	case "rdp":
		info := populatePortsInfo([]int{3389}, "TCP + UDP", wikiString + getWikiName(service))
		return info

	case "ldap":
		info := populatePortsInfo([]int{389}, "TCP + UDP", wikiString + getWikiName(service))
		return info

	case "bgp":
		info := populatePortsInfo([]int{179}, "TCP", wikiString + getWikiName(service))
		return info
	
	// Most of the common ports can be retrieved via the in-built net package
	default:
		transportProtocol := "tcp"
		portNum, err := net.LookupPort(transportProtocol, service)
		if err != nil {
			log.Fatal(err)
		}

		info := populatePortsInfo([]int{portNum}, strings.ToUpper(transportProtocol), wikiString + getWikiName(service))
		return info

	}
}

func ipv4PrivateAddressRange() {
	fmt.Println("The IPv4 private address spaces are:")
	fmt.Printf("\t10.0.0.0 - 10.255.255.255\n\t172.16.0.0 - 172.31.255.255\n\t192.168.0.0 - 192.168.255.255\n")
}

func main() {
	myIP := "192.168.1.5/20"
	ip, ipnet := parseIPInfo(myIP)
	fmt.Println("IP",ip, ipnet)
	testProto := getCommonPorts("ldap")
	printPortInfo(testProto)
}