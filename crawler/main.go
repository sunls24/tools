package main

import (
	"crawler/options"
	"crawler/utils"
	"crawler/yushugu"
	"flag"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

const (
	defaultUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
)

func main() {
	var opts = options.ParseOptions()
	if opts.Help {
		flag.PrintDefaults()
		return
	}
	fmt.Printf("%+v\n", opts)
	if len(opts.Output) == 0 {
		panic("output path is empty")
	}

	ff := os.O_CREATE | os.O_WRONLY
	if opts.Append {
		ff = ff | os.O_APPEND
	}
	file, err := os.OpenFile(opts.Output, ff, 0o644)
	utils.Check(err)
	defer file.Close()

	var ua = defaultUA

	c := colly.NewCollector(
		colly.UserAgent(ua), // 设置UA
		colly.Async(true),   // 异步
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.StatusCode, err)
	})

	switch opts.Provider {
	case options.Yushugu:
		err = yushugu.Visit(file, c, opts)
	default:
		panic(fmt.Sprintf("no provider: %s", opts.Provider))
	}

	utils.Check(err)

	c.Wait()
}
