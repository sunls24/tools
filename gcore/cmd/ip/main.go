package main

import (
	"flag"
	"gcore/ip"
)

var ipv6 bool

func init() {
	flag.BoolVar(&ipv6, "v6", false, "ipv6")
	flag.Parse()
}

func main() {
	ipList := ip.GetGcoreIp()
	if ipv6 {
		ip.PrintIp(ipList.AddressesV6)
	} else {
		ip.PrintIp(ipList.Addresses)
	}
}
