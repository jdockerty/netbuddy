package main

import (
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"net"
	"flag"
	"log"
	"os"
	"strings"
)

type SubnetInfo struct {
	networkAddress      net.IP
	broadcastAddress    net.IP
	firstUsuableAddress net.IP
	lastUsuableAddress  net.IP
	totalAddressCount   uint64
}

type PortsInfo struct {
	commonPortNumbers      []int
	transportLayerProtocol string
	extraInfoLink          string
}

func parseIPInfo(ipString string) (net.IP, *net.IPNet) {
	ip, ipnet, err := net.ParseCIDR(ipString)
	if err != nil {
		log.Fatal("Error parsing IP and CIDR notation, ensure it looks something like '192.168.3.1/24'")
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

func subnetIterations(ipnet *net.IPNet, iterations int) {
	_, lastAddress := getNetworkAndBroadcast(ipnet)
	prefixBits, _ := ipnet.Mask.Size()

	for i := 1; i < iterations+1; i++ {
		nextNetworkAddress := cidr.Inc(lastAddress)
		nextPrefixString := fmt.Sprint(nextNetworkAddress, "/", prefixBits)
		fmt.Printf("[%d] Next subnet: %s\n", i, nextPrefixString)
		_, nextIPNet := parseIPInfo(nextPrefixString)
		_, lastAddress = getNetworkAndBroadcast(nextIPNet)

		nextSubnetInfo := getSubnetInfo(nextIPNet)
		printSubnetInfo(nextSubnetInfo)

	}

}

func printSubnetInfo(subnet SubnetInfo) {
	fmt.Printf("Network: %s\nFirst assignable: %s\nLast assignable: %s\nBroadcast: %s\n",
		subnet.networkAddress, subnet.firstUsuableAddress, subnet.lastUsuableAddress, subnet.broadcastAddress)
}

func getSubnetInfo(ipnet *net.IPNet) SubnetInfo {
	networkAddr, broadcastAddr := getNetworkAndBroadcast(ipnet)
	firstUsuableAddr := cidr.Inc(networkAddr)
	lastUsuableAddr := cidr.Dec(broadcastAddr)
	totalAddresses := getAddressCount(ipnet)

	subnetInfoResponse := SubnetInfo{
		networkAddress:      networkAddr,
		broadcastAddress:    broadcastAddr,
		firstUsuableAddress: firstUsuableAddr,
		lastUsuableAddress:  lastUsuableAddr,
		totalAddressCount:   totalAddresses,
	}

	return subnetInfoResponse
}

func populatePortsInfo(portNums []int, transportProto, link string) PortsInfo {
	portsInfoResponse := PortsInfo{
		commonPortNumbers:      portNums,
		transportLayerProtocol: transportProto,
		extraInfoLink:          link,
	}
	return portsInfoResponse
}

func getWikiName(service string) string {

	abbreviationsToWikiName := map[string]string{
		"dns":    "Domain_Name_System",
		"dhcp":   "Dynamic_Host_Configuration_Protocol",
		"rdp":    "Remote_Desktop_Protocol",
		"smtp":   "Simple_Network_Management_Protocol",
		"ssh":    "Secure_Shell",
		"telnet": "telnet",
		"ftp":    "File_Transfer_Protocol",
		"http":   "Hypertext_Transfer_Protocol",
		"https":  "HTTPS",
		"imap":   "Internet_Message_Access_Protocol",
		"pop3":   "Post_Office_Protocol",
		"ldap":   "Lightweight_Directory_Access_Protocol",
		"bgp":    "Border_Gateway_Protocol",
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
		info := populatePortsInfo([]int{53}, "UDP", wikiString+getWikiName(service))
		return info

	case "dhcp":
		info := populatePortsInfo([]int{67, 68}, "UDP", wikiString+getWikiName(service))
		return info

	case "rdp":
		info := populatePortsInfo([]int{3389}, "TCP + UDP", wikiString+getWikiName(service))
		return info

	case "ldap":
		info := populatePortsInfo([]int{389}, "TCP + UDP", wikiString+getWikiName(service))
		return info

	case "bgp":
		info := populatePortsInfo([]int{179}, "TCP", wikiString+getWikiName(service))
		return info

	// Most of the common ports can be retrieved via the in-built net package
	default:
		transportProtocol := "tcp"
		portNum, err := net.LookupPort(transportProtocol, service)
		if err != nil {
			fmt.Printf("Unsupported service lookup: %s\n", service)
			log.Fatal(err)
		}

		info := populatePortsInfo([]int{portNum}, strings.ToUpper(transportProtocol), wikiString+getWikiName(service))
		return info

	}
}

func ipv4PrivateAddressRange() {
	fmt.Println("The RFC 1918 IPv4 private address spaces are:")
	fmt.Printf("\t10.0.0.0 - 10.255.255.255\n\t172.16.0.0 - 172.31.255.255\n\t192.168.0.0 - 192.168.255.255\n")
}

func subnetCmdHelp() {
	fmt.Println("Usage: netbuddy subnet <arg> <input>")
	fmt.Println("Args:\n\t-display: Shows various information about a particular IP and CIDR, e.g. 192.168.4.20/19")
	fmt.Println("\t-count: Show the total number of addresses in the provided network.")
	fmt.Println("\t-iterate: Show the next X iterations of a particular prefix to the network.")
	fmt.Println("\nExamples: \n\t netbuddy subnet -count 172.31.5.9/19\n\t netbuddy subnet -iterate 2 192.168.0.0/24")
}
func showCmdHelp() {
	fmt.Println("Usage: netbuddy show <option> <input>")
	fmt.Println("Options:\n\tipv4range - Show RFC 1918 IPv4 address range. \tNote: This does not take an input.")
	fmt.Println("\tmac - Show MAC address of an interface.")
	fmt.Println("\tservice - Shows port and information for a particular service e.g. SSH")
	fmt.Println("\nExamples: \n\t netbuddy show service ssh\n\t netbuddy show ipv4range")
}

func main() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	subnetCmd := flag.NewFlagSet("subnet", flag.ExitOnError)
	subnetDisplay := subnetCmd.String("display", "", "Displays the various addresses within a given subnet.")
	subnetIterate := subnetCmd.Int("iterate", 0, "Iterates over and displays the next X networks for a prefix.")
	subnetAddressCount := subnetCmd.String("count", "", "Displays the total available addresses for a given network.")

	switch os.Args[1] {
	case "show":
		showCmd.Parse(os.Args[2:])
		switch os.Args[2] {
		case "ipv4range":
			ipv4PrivateAddressRange()
		case "interfaces":
			fmt.Println("TO DO")
		case "service":
			portInfo := getCommonPorts(os.Args[3])
			printPortInfo(portInfo)
		case "help":
			showCmdHelp()
		default:
			showCmdHelp()
		}

	case "subnet":
		subnetCmd.Parse(os.Args[2:])

		if os.Args[2] == "help" {
			subnetCmdHelp()
		}
		if len(*subnetDisplay) != 0 {
			_, ipnet := parseIPInfo(*subnetDisplay)
			subnetInfo := getSubnetInfo(ipnet)
			printSubnetInfo(subnetInfo)
		}

		if *subnetIterate > 0 {
			_, ipnet := parseIPInfo(os.Args[4])
			subnetIterations(ipnet, *subnetIterate)
		}

		if len(*subnetAddressCount) != 0 {
			_, ipnet := parseIPInfo(*subnetAddressCount)
			subnetInfo := getSubnetInfo(ipnet)
			fmt.Printf("There are %d total available addresses in this network.\n", subnetInfo.totalAddressCount)
		}

	}

}
