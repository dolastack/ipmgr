package ipconv

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Ipfacts struct {
	ipAddress      []int
	gateway        []int
	brdcastAddress []int
	netmask        []int
	mask           int
}

func rightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func ParsInput(cidr string) *Ipfacts {
	/** parse **/
	var newip *Ipfacts
	ip := make([]int, 4)
	var netmask int
	pattern := regexp.MustCompile(`(?P<ip>\d{1,3}(?:\.\d{1,3}){3})/(?P<mask>[0-9]{1,2})$`)
	if pattern.MatchString(cidr) {
		groups := pattern.FindAllStringSubmatch(cidr, 2)

		ips := strings.Split(groups[0][1], ".")
		netmask, _ = strconv.Atoi(groups[0][2])
		for i, _ := range ips {
			ip[i], _ = strconv.Atoi(ips[i])
		}
		newip = NewIpfacts()
		newip.ipAddress = ip
		//validate ip
		if !validateIP(ip) {
			println("invalid IP")
			os.Exit(1)
		}
		//validate mask
		if !validateMask(netmask) {
			println("invalid netmask")
			os.Exit(1)
		}
	} else {
		println("Invalid input, Input CIDR format a.b.c.d/m")
		os.Exit(1)
	}
	return newip
}

func validateMask(mask int) bool {
	return mask > 1 && mask <= 32
}
func validateIP(ip []int) bool {
	test := false
	for _, oc := range ip {
		if oc >= 0 && oc <= 255 {
			test = true
		} else {
			test = false
			break
		}
	}
	return test
}

func MaskToBinary(netmask int) string {
	mask := ""
	for i := 0; i < netmask; i++ {
		mask += "1"
	}
	//println(mask)
	mask = rightPad2Len(mask, "0", 32)
	return mask
}

func IpToBinary(ip []int) string {
	var octectsB [4]string
	ipbinary := ""
	if validateIP(ip) {
		//octets := strings.Split(ip, ".")
		for i, octet := range ip {
			octectsB[i] = strconv.FormatInt(int64(octet), 2)

			octectsB[i] = leftPad2Len(octectsB[i], "0", 8)
			ipbinary += octectsB[i]
		}

		//fmt.Println(octectsB)
		//println(ipbinary)
		return ipbinary
	} else {
		fmt.Println(ip, "is not a valid IP address")
		return ""
	}
}

func binToDec(bin string) int {
	acc := 0
	lent := len(bin)
	for i := 0; i < lent; i++ {
		v := int(bin[i]) - '0' //rune conversion
		acc += v * int(math.Pow(2, float64(lent-i-1)))
	}
	return acc
}
func BinaryToIp(ip string) []string {
	octets := make([]string, 4)
	for i, j := 0, 0; i < 32; i += 8 {
		//println(ip[i : i+8])
		octet := binToDec(ip[i : i+8])
		octets[j] = strconv.Itoa(octet)
		j++
	}
	//fmt.Println(octets)
	return octets
}

func GetSubnet(ip []string, mask []string) []string {
	subnet := make([]string, 4)
	for i, j := 0, 0; i < 4; i++ {
		y, _ := strconv.Atoi(ip[i])
		z, _ := strconv.Atoi(mask[i])
		subnet[j] = strconv.Itoa(y | z)
		j++
	}
	//fmt.Println(subnet)
	return subnet
}

func GetSubnetFacts(ip *Ipfacts) []string {
	var rangeips func(ra []int)
	facts := make([]string, 4)
	var octet int
	switch {
	case ip.mask < 8:
		octet = 0
		rangeips = func(ra []int) {
			ip.gateway = []int{ra[0], 0, 0, 1}
			ip.brdcastAddress = []int{ra[1] - 1, 255, 255, 255}
		}
	case ip.mask < 16:
		octet = 1
		rangeips = func(ra []int) {

			ip.gateway = []int{ip.ipAddress[0], ra[0], 0, 1}
			ip.brdcastAddress = []int{ip.ipAddress[0], ra[1] - 1, 255, 255}

		}
	case ip.mask < 24:
		octet = 2
		rangeips = func(ra []int) {

			ip.gateway = []int{ip.ipAddress[0], ip.ipAddress[1], ra[0], 1}
			ip.brdcastAddress = []int{ip.ipAddress[0], ip.ipAddress[1], ra[1] - 1, 255}
		}
	case ip.mask < 32:
		octet = 3
		rangeips = func(ra []int) {

			ip.gateway = []int{ip.ipAddress[0], ip.ipAddress[1], ip.ipAddress[2], ra[0] + 1}
			ip.brdcastAddress = []int{ip.ipAddress[0], ip.ipAddress[1], ip.ipAddress[2], ra[1] - 1}

		}
	}

	multple := int(math.Pow(2, float64(8-(ip.mask%8))))
	//fmt.Println(multple)
	size := 256/multple + 1
	multples := make([]int, size)
	for i := 0; i < size; i++ {
		multples[i] += i * multple
	}
	//fmt.Println(multples)
	rangip := func() []int {
		out := make([]int, 2)
		for i, it := range multples {
			if it > ip.ipAddress[octet] {
				out[0] = multples[i-1]
				out[1] = multples[i]
				return out
			}
		}
		return out
	}()

	rangeips(rangip)

	return facts
}

func DisplayFacts(ip *Ipfacts) {
	fmt.Println("IP Address : ", ip.ipAddress)
	fmt.Println("IP Netmask : ", ip.netmask)
	fmt.Println("IP Mask : ", ip.mask)
	fmt.Println("Gateway : ", ip.gateway)
	fmt.Println("Braodcast Address: ", ip.brdcastAddress)
}

func NewIpfacts() *Ipfacts {
	return &Ipfacts{
		ipAddress:      []int{0, 0, 0, 0},
		gateway:        []int{0, 0, 0, 0},
		brdcastAddress: []int{0, 0, 0, 0},
		netmask:        []int{0, 0, 0, 0},
		mask:           0,
	}
}
