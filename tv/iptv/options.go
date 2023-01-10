package iptv

import "flag"

type Options struct {
	Source    string // 源地址，文件路径或者URL
	Available string // 已经检查过可用性的源，仅需要排序时使用
	Output    string // 输出路径
	EpgXml    string // epg 地址
	EpgOutput string // epg 输出地址
	Workdir   string // 工作目录
	Parallel  int    // iptv-checker 线程数量
	Help      bool
}

func ParseOptions() Options {
	var opts Options
	flag.StringVar(&opts.Source, "s", "", "Source path or url")
	flag.StringVar(&opts.Available, "a", "", "Available path or url")
	flag.StringVar(&opts.Output, "o", "online.m3u", "Output file path")
	flag.StringVar(&opts.EpgOutput, "eo", "online.xml", "Epg xml file path")
	flag.StringVar(&opts.EpgXml, "e", "", "Epg xml path or url")
	flag.StringVar(&opts.Workdir, "w", "workdir", "workdir")
	flag.IntVar(&opts.Parallel, "p", 100, "Parallel")
	flag.BoolVar(&opts.Help, "h", false, "Help")
	flag.Parse()

	if opts.Help {
		return opts
	}

	if len(opts.Source) == 0 {
		panic("Souce is empty")
	}
	if len(opts.Output) == 0 {
		panic("Output is empty")
	}
	if opts.Parallel <= 0 {
		opts.Parallel = 1
	}
	return opts
}
