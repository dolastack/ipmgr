package main

import (
	"fmt"
	"ipconv"
	"os/exec"
	"regexp"
	"strings"
)

type iface struct {
	name string
	ips  []ipconv.Ipfacts
}

func main() {
	ifaces := getIfname()
	for _, it := range ifaces {
		if !(it == "lo" || it == "") {
			println(it)

			cmd := fmt.Sprintf("ip a show %s | egrep  inet ", it)

			inets := runCmd(cmd)
			inetsA := strings.Split(inets, "\n")
			for _, in := range inetsA {
				getIp4(in)
			}
		}
	}

}

func runCmd(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err == nil {
		return string(out)
	}
	return ""
}

func getIfname() []ifaces {
	//get the active interfaces on the system
	var ifaces []iface
	var names []string
	cmd := `netstat -i |  cut -f 1 -d ' '`
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err == nil {
		outs := strings.Split(string(out), "\n")
		names = outs[2:]
		fmt.Println(ifaces)

		ifaces = iface{
			name: ifaces,
			ips:  []ipconv.Newipfacts(),
		}
	}
	return ifaces
}

func getIp4(inet string) {
	pattern := regexp.MustCompile(`\b(?P<ip>\d{1,3}(?:\.\d{1,3}){3})/(?P<mask>[0-9]{1,2})\b`)
	if pattern.MatchString(inet) {
		groups := pattern.FindAllStringSubmatch(inet, 2)
		fmt.Println(groups[0][0])
	}
}

func NewIfaces(temp []ipconv.Ipfacts) *iface {
	return &iface{
		name: "",
		ips:  temp,
	}
}
