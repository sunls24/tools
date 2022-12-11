package yushugu

import (
	"crawler/options"
	"fmt"
	"os"
	"time"

	"github.com/gocolly/colly"
)

const (
	baseUrl         = "https://m.yushugu.com/"
	defaultInterval = time.Millisecond * 500
)

func Visit(file *os.File, c *colly.Collector, opts options.Options) error {
	if len(opts.FirstUrl) == 0 {
		panic("first url is empty")
	}

	var interval = opts.Interval
	if interval == 0 {
		interval = defaultInterval
		fmt.Println("use default interval:", interval)
	}
	c.OnHTML("section>h3", func(e *colly.HTMLElement) {
		fmt.Println(">", e.Text)
		if opts.WithChapter {
			file.WriteString(e.Text)
			file.WriteString("\n")
		}
	})

	c.OnHTML("section>p", func(e *colly.HTMLElement) {
		if len(e.Text) == 0 {
			return
		}
		file.WriteString(e.Text)
		file.WriteString("\n")
	})

	c.OnHTML(".btn-next", func(e *colly.HTMLElement) {
		if e.Text == "目录" {
			return
		}
		if e.Text == "没有了" {
			fmt.Println("✅没有了")
			return
		}
		<-time.After(interval)
		_ = c.Visit(baseUrl + e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	return c.Visit(baseUrl + opts.FirstUrl)
}
