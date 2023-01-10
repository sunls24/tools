package main

import (
	"flag"
	"fmt"
	"tv/iptv"
)

func main() {
	opts := iptv.ParseOptions()
	if opts.Help {
		flag.PrintDefaults()
		return
	}
	fmt.Printf("opts: %+v\n", opts)
	iptv.Checker(opts)
}
