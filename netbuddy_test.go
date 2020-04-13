package main

import (
	"testing"
	"fmt"
)


func TestGetAddressCount(t *testing.T) {
	_, ipnet := parseIPInfo("192.168.1.1/24")
	addressCount := getAddressCount(ipnet)
	fmt.Println("Total address count for", ipnet, "is", addressCount)
	if addressCount != 256 {
		t.Errorf("/24 CIDR returned %d, required 256.", addressCount)
	}
}



