package main

import (
	"flag"
	"io/ioutil"
	"strings"

	"github.com/yene/domain-check"
)

var domain = flag.String("domain", "yannickweiss.com", "Domain to check")
var hostsFile = flag.String("hosts", "", "The path to the file containing a list of hosts to check.")
var warnDays = flag.Int("days", 0, "Warn if the certificate will expire within this many days.")

func main() {
	flag.Parse()
	if len(*hostsFile) == 0 {
		flag.Usage()
		return
	}

	for _, host := range hostsFromFile(*hostsFile) {
		checker.CheckDomain(host)
	}
}

func hostsFromFile(hostFile string) (hosts []string) {
	fileContents, err := ioutil.ReadFile(*hostsFile)
	if err != nil {
		return hosts
	}
	lines := strings.Split(string(fileContents), "\n")
	for _, line := range lines {
		host := strings.TrimSpace(line)
		if len(host) == 0 || host[0] == '#' {
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts
}
