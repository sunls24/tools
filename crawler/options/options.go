package options

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// provider
const (
	Yushugu = "yushugu"
)

type Options struct {
	Provider    string
	FirstUrl    string
	WithChapter bool
	Append      bool
	Interval    time.Duration
	Output      string
	Help        bool
}

func ParseOptions() Options {
	var opts Options
	flag.StringVar(&opts.Provider, "p", Yushugu, "provider")
	flag.StringVar(&opts.FirstUrl, "fu", "", "first url")
	flag.BoolVar(&opts.WithChapter, "wc", false, "with chapter")
	flag.BoolVar(&opts.Append, "ap", false, "append file")
	var interval string
	flag.StringVar(&interval, "i", "", "interval")
	flag.StringVar(&opts.Output, "o", "", "output path")
	flag.BoolVar(&opts.Help, "h", false, "help")
	flag.Parse()
	if len(interval) != 0 {
		opts.Interval, _ = time.ParseDuration(interval)
	}
	if len(opts.Output) != 0 && !strings.HasSuffix(opts.Output, "txt") {
		opts.Output = fmt.Sprintf("%s.txt", opts.Output)
	}
	return opts
}
