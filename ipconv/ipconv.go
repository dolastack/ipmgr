package ipconv

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type cidr string

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
	var mask int
	pattern := regexp.MustCompile(`(?P<ip>\d{1,3}(?:\.\d{1,3}){3})/(?P<mask>[0-9]{1,2})$`)
	if pattern.MatchString(string(cidr)) {
		groups := pattern.FindAllStringSubmatch(string(cidr), 2)

		ips := strings.Split(groups[0][1], ".")
		mask, _ = strconv.Atoi(groups[0][2])

		//validate ip
		if !validateIP(ip) {
			println("invalid IP")
			os.Exit(1)
		}
		//validate mask
		if !validateMask(mask) {
			println("invalid netmask")
			os.Exit(1)
		}

		for i, _ := range ips {
			ip[i], _ = strconv.Atoi(ips[i])
		}
		newip = NewIpfacts(ip, mask)

		newip.netmask = binaryToIp(newip.maskToBinary())

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

func (ipf *Ipfacts) maskToBinary() string {
	//convert mask (int) to a string of ones and zeros
	mask := ""
	for i := 0; i < ipf.mask; i++ {
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
func binaryToIp(ip string) []int {
	octets := make([]int, 4)
	for i, j := 0, 0; i < 32; i += 8 {
		//println(ip[i : i+8])
		octets[j] = binToDec(ip[i : i+8])
		j++
	}
	//fmt.Println(octets)
	return octets
}

func GetSubnet(ip []int, mask []string) []string {
	subnet := make([]string, 4)
	for i, j := 0, 0; i < 4; i++ {
		y := ip[i]
		z, _ := strconv.Atoi(mask[i])
		subnet[j] = strconv.Itoa(y | z)
		j++
	}
	//fmt.Println(subnet)
	return subnet
}

func (ip *Ipfacts) GetSubnetFacts() []string {
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
	//fmt.Println(ip.mask)
	multple := int(math.Pow(2, float64(8-(ip.mask%8))))
	//fmt.Println(multple)
	size := 256/multple + 1
	multples := make([]int, size)
	for i := 0; i < size; i++ {
		//multiples of range of subnet
		multples[i] += i * multple
	}
	//fmt.Println(multples)
	rangip := func() []int {
		out := make([]int, 2)
		for i, it := range multples {
			if it > ip.ipAddress[octet] {
				out[0], out[1] = multples[i-1], multples[i]
				return out
			}
		}
		return out
	}()

	rangeips(rangip)

	return facts
}

func (ip *Ipfacts) DisplayFacts() {
	fmt.Println("IP Address : ", ip.ipAddress)
	fmt.Println("IP Netmask : ", ip.netmask)
	fmt.Println("IP Mask : ", ip.mask)
	fmt.Println("Gateway : ", ip.gateway)
	fmt.Println("Braodcast Address: ", ip.brdcastAddress)
}

func NewIpfacts(ipad []int, mask int) *Ipfacts {
	return &Ipfacts{
		ipAddress:      ipad,
		gateway:        []int{0, 0, 0, 0},
		brdcastAddress: []int{0, 0, 0, 0},
		netmask:        []int{0, 0, 0, 0},
		mask:           mask,
	}
}
