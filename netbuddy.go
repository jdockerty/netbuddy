package main

import (
	"strconv"
	"flag"
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"log"
	"net"
	"os"
	"strings"
)

type subnetData struct {
	networkAddress      net.IP
	broadcastAddress    net.IP
	firstUsuableAddress net.IP
	lastUsuableAddress  net.IP
	totalAddressCount   uint64
}

type portsData struct {
	commonPortNumbers      []int
	transportLayerProtocol string
	extraInfoLink          string
	err                    int
}

type interfaceData struct {
	name    string
	ipAddr  string
	macAddr string
	flags   string
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
	prefixBits, _ := ipnet.Mask.Size()

	for i := 1; i < iterations+1; i++ {
		nextPrefixString, _ := cidr.NextSubnet(ipnet, prefixBits)
		fmt.Printf("\n[%d] Next subnet: %s\n", i, nextPrefixString)

		_, nextIPNet := parseIPInfo(nextPrefixString.String())
		nextsubnetData := getSubnetData(nextIPNet)

		printSubnetData(nextsubnetData)
		ipnet = nextPrefixString
	}

}

func printSubnetData(subnet subnetData) {
	fmt.Printf("Network: %s\nFirst assignable: %s\nLast assignable: %s\nBroadcast: %s\n",
		subnet.networkAddress, subnet.firstUsuableAddress, subnet.lastUsuableAddress, subnet.broadcastAddress)
}

func getSubnetData(ipnet *net.IPNet) subnetData {
	networkAddr, broadcastAddr := getNetworkAndBroadcast(ipnet)
	firstUsuableAddr := cidr.Inc(networkAddr)
	lastUsuableAddr := cidr.Dec(broadcastAddr)
	totalAddresses := getAddressCount(ipnet)

	subnetDataResponse := subnetData{
		networkAddress:      networkAddr,
		broadcastAddress:    broadcastAddr,
		firstUsuableAddress: firstUsuableAddr,
		lastUsuableAddress:  lastUsuableAddr,
		totalAddressCount:   totalAddresses,
	}

	return subnetDataResponse
}

func populatePortsData(portNums []int, transportProto, link string, errorNum int) portsData {
	portsDataResponse := portsData{
		commonPortNumbers:      portNums,
		transportLayerProtocol: transportProto,
		extraInfoLink:          link,
		err:                    errorNum,
	}
	return portsDataResponse
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

func printPortData(portInfo portsData) {
	if portInfo.err != 1 {
		fmt.Printf("Port Numbers: %d\nTransport Protocol(s): %s\nFor more information on this protocol visit %s\n",
			portInfo.commonPortNumbers, portInfo.transportLayerProtocol, portInfo.extraInfoLink)
	}

}

func getCommonPorts(service string) portsData {

	wikiString := "https://en.wikipedia.org/wiki/"
	service = strings.ToLower(service)

	switch service {
	case "dns":
		info := populatePortsData([]int{53}, "UDP", wikiString+getWikiName(service), 0)
		return info

	case "dhcp":
		info := populatePortsData([]int{67, 68}, "UDP", wikiString+getWikiName(service), 0)
		return info

	case "rdp":
		info := populatePortsData([]int{3389}, "TCP + UDP", wikiString+getWikiName(service), 0)
		return info

	case "ldap":
		info := populatePortsData([]int{389}, "TCP + UDP", wikiString+getWikiName(service), 0)
		return info

	case "bgp":
		info := populatePortsData([]int{179}, "TCP", wikiString+getWikiName(service), 0)
		return info

	// Most of the common ports can be retrieved via the in-built net package
	default:
		transportProtocol := "tcp"
		portNum, err := net.LookupPort(transportProtocol, service)
		if err != nil {
			fmt.Printf("Unsupported service lookup: %s\n", service)
			info := populatePortsData([]int{}, "", "", 1)
			return info
		}

		info := populatePortsData([]int{portNum}, strings.ToUpper(transportProtocol), wikiString+getWikiName(service), 0)
		return info

	}
}

func convertToCIDRNotation(mask string) int {
	maskAsSlice := strings.Split(mask, ".")
	var byteMaskSlice []byte
	for _, value := range maskAsSlice {
		stringToInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Error parsing subnet mask, ensure it is in the form X.X.X.X")
		}
		byteMaskSlice = append(byteMaskSlice, byte(stringToInt))
	}
	maskAsIPMaskType := net.IPv4Mask(byteMaskSlice[0], byteMaskSlice[1], byteMaskSlice[2], byteMaskSlice[3]) // .IPv4Mask expects form (a, b, c, d byte)
	sizeOfMask, _ := maskAsIPMaskType.Size()
	fmt.Printf("%s is equivalent to the CIDR notation of /%d\n", mask, sizeOfMask)
	return sizeOfMask
}

func showinterfaceData() bool {
	interfaceSlice, _ := net.Interfaces()
	for _, interfaceData := range interfaceSlice {
		currentAddr, err := interfaceData.Addrs()
		if err != nil {
			fmt.Println("Could not get interface address details.")
		}
		fmt.Printf("Interface name: %s\nAssociated IP Addresses: %s\nMAC Address: %s\nOther Info: %s\n\n",
			interfaceData.Name,
			currentAddr,
			interfaceData.HardwareAddr,
			interfaceData.Flags)
	}
	return true
}

func ipv4PrivateAddressRange() bool {
	fmt.Println("The RFC 1918 IPv4 private address spaces are:")
	fmt.Printf("\t10.0.0.0 - 10.255.255.255\n\t172.16.0.0 - 172.31.255.255\n\t192.168.0.0 - 192.168.255.255\n")
	return true
}

func subnetCmdHelp() bool {
	fmt.Println("Usage: netbuddy subnet <arg> <input>")
	fmt.Println("Args:\n\t-display: Shows various information about a particular IP and CIDR, e.g. 192.168.4.20/19")
	fmt.Println("\t-count: Show the total number of addresses in the provided network.")
	fmt.Println("\t-iterate: Show the next X iterations of a particular prefix to the network.")
	fmt.Println("\nExamples: \n\t netbuddy subnet -count 172.31.5.9/19\n\t netbuddy subnet -iterate 2 192.168.0.0/24")
	return true
}
func showCmdHelp() bool {
	fmt.Println("Usage: netbuddy show <option> <input>")
	fmt.Println("Options:\n\tipv4range - Show RFC 1918 IPv4 address range. \tNote: This does not take an input.")
	fmt.Println("\tinterfaces - Show interface information on this machine.")
	fmt.Println("\tservice - Shows port and information for a particular service e.g. SSH")
	fmt.Println("\nExamples: \n\t netbuddy show service ssh\n\t netbuddy show ipv4range")
	return true
}

func main() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	subnetCmd := flag.NewFlagSet("subnet", flag.ExitOnError)
	subnetDisplay := subnetCmd.String("display", "", "Displays the various addresses within a given subnet.")
	subnetIterate := subnetCmd.Int("iterate", 0, "Iterates over and displays the next X networks for a prefix.")
	subnetAddressCount := subnetCmd.String("count", "", "Displays the total available addresses for a given network.")
	subnetDecToCIDR := subnetCmd.String("tocidr", "", "Convert the dotted decimal notation of an IPv4 subnet mask to the equivalent CIDR.")

	switch os.Args[1] {
	case "show":
		showCmd.Parse(os.Args[2:])
		switch os.Args[2] {
		case "ipv4range":
			ipv4PrivateAddressRange()
		case "interfaces":
			showinterfaceData()
		case "service":
			portInfo := getCommonPorts(os.Args[3])
			printPortData(portInfo)
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
			subnetData := getSubnetData(ipnet)
			printSubnetData(subnetData)
		}

		if *subnetIterate > 0 {
			_, ipnet := parseIPInfo(os.Args[4])
			subnetIterations(ipnet, *subnetIterate)
		}

		if len(*subnetAddressCount) != 0 {
			_, ipnet := parseIPInfo(*subnetAddressCount)
			subnetData := getSubnetData(ipnet)
			fmt.Printf("There are %d total available addresses in this network.\n", subnetData.totalAddressCount)
		}

		if len(*subnetDecToCIDR) != 0 {
			convertToCIDRNotation(*subnetDecToCIDR)
		}
	default:
		fmt.Printf("The currently supported commands are: \n- show\n- subnet\n Use 'netbuddy <command> help' for more information.\n")
	}

}
