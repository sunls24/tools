package main

import (
	"flag"
	"gcore/ping"
)

var pingCount, limit int
var countryCode string

func init() {
	flag.IntVar(&pingCount, "c", 100, "ping count")
	flag.IntVar(&limit, "l", 100, "limit count")
	flag.StringVar(&countryCode, "cc", "", "country code")
	flag.Parse()
}

func main() {
	ping.PingGcoreIp(pingCount, limit, countryCode)
}
