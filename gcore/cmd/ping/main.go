package main

import (
	"flag"
	"gcore/ping"
)

var pingCount int

func init() {
	flag.IntVar(&pingCount, "c", 100, "ping count")
	flag.Parse()
}

func main() {
	ping.PingGcoreIp(pingCount)
}
