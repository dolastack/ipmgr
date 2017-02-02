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

	//ips := ipconv.BinaryToIp(ipconv.IpToBinary(ip))
	//mask := ipconv.BinaryToIp(ipconv.MaskToBinary(netmask))
	//	ipconv.GetSubnet(ips, mask)
	ipconv.GetSubnetFacts(ip)
	ipconv.DisplayFacts(ip)
}
