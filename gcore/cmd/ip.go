package main

import "flag"

var ipv6 bool

func init() {
	flag.BoolVar(&ipv6, "v6", false, "ipv6")
	flag.Parse()
}

func main() {
}
