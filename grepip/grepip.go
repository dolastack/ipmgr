package main

import (
	"fmt"
	"ipmgr/ipconv"
	"os/exec"
	"regexp"
	"strings"
)

type ifname string
type cidr string

type iface struct {
	name ifname
	ips  []*ipconv.Ipfacts
}

func main() {
	ifaces := getIfname()
	for _, temp := range ifaces {
		//	println(ifc)
		ifc := getDetails(temp)
		fmt.Println(ifc.name)
		for _, i := range ifc.ips {
			i.DisplayFacts()
		}

	}
}

func runCmd(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err == nil {
		return string(out)
	} else {
		return ""
	}
}

func getIfname() []ifname {
	//get the active interfaces on the system
	var ifaces []ifname
	cmd := `netstat -i |  cut -f 1 -d ' '`
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err == nil {
		outs := strings.Split(string(out), "\n")
		for _, it := range outs[2:] {
			if !(it == "" || it == "lo") {
				ifaces = append(ifaces, ifname(it))
			}
		}
		//	fmt.Println(ifaces)
	}
	return ifaces
}

func getCidr(inet string) cidr {
	var thecidr cidr
	pattern := regexp.MustCompile(`\b(?P<ip>\d{1,3}(?:\.\d{1,3}){3})/(?P<mask>[0-9]{1,2})\b`)
	if pattern.MatchString(inet) {
		groups := pattern.FindAllStringSubmatch(inet, 2)
		thecidr = cidr(groups[0][0])
		//fmt.Println(thecidr)
	}
	return thecidr
}

func getDetails(ifc ifname) *iface {
	//	var ifc
	var ipfacts []*ipconv.Ipfacts
	cmd := fmt.Sprintf(`ip a show %s | grep -P 'inet\s'`, ifc)

	inets := runCmd(cmd)
	inetsA := strings.Split(inets, "\n")
	//fmt.Println(inetsA)
	for _, in := range inetsA {
		if in != "" {
			//	fmt.Println(getCidr(in))
			ipf := ipconv.ParsInput(string(getCidr(in)))
			ipf.GetSubnetFacts()
			ipfacts = append(ipfacts, ipf)
		}
	}
	return NewIfaces(ifc, ipfacts)
}

func NewIfaces(ifc ifname, temp []*ipconv.Ipfacts) *iface {
	//constructor for iface struct
	return &iface{
		name: ifc,
		ips:  temp,
	}
}
