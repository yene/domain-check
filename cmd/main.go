package main

import (
	"flag"

	"github.com/yene/domain-check"
)

var domain = flag.String("domain", "yannickweiss.com", "Domain to check")

func main() {
	flag.Parse()
	checker.CheckDomain(*domain)
}
