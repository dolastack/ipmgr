package main

import (
	"fmt"
	"ipmgr/ipconv"
)

func main() {
	cidr := ""

	fmt.Print("Enter IP address : ")
	fmt.Scanf("%s", &cidr)
	ip := ipconv.ParsInput(cidr)

	ip.GetSubnetFacts()
	ip.DisplayFacts()
}
