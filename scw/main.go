package main

import (
	"flag"
	"fmt"
	"log"
	"scw/command"
	"time"
)

var (
	StartCount     int
	StartInterval  string
	StatusInterval string
	ProjectId      string

	PrintCommand bool
	DeleteAtEnd  bool
	CreateAtEnd  bool
	Help         bool
)

func init() {
	flag.IntVar(&StartCount, "sc", 10, "开机次数")
	flag.StringVar(&StartInterval, "si", "10m", "开机间隔")
	flag.StringVar(&StatusInterval, "ssi", "1m", "检测状态间隔")
	flag.StringVar(&ProjectId, "pid", "", "project id")
	flag.BoolVar(&PrintCommand, "pc", false, "是否输出执行的命令")
	flag.BoolVar(&DeleteAtEnd, "de", false, "循环结束后是否删除机器")
	flag.BoolVar(&CreateAtEnd, "ce", false, "循环结束后是否删除并重建机器，需要和de同时开启")
	flag.BoolVar(&Help, "h", false, "命令说明")
	flag.Parse()
}

func main() {
	if Help || ProjectId == "" {
		fmt.Println("Required for -pid args")
		flag.PrintDefaults()
		return
	}

	si, err := time.ParseDuration(StartInterval)
	if err != nil {
		log.Panicln(err)
	}
	ssi, err := time.ParseDuration(StatusInterval)
	if err != nil {
		log.Panicln(err)
	}
	opts := command.Options{
		StartCount:     StartCount,
		StartInterval:  si,
		StatusInterval: ssi,
		PrintCommand:   PrintCommand,
		ProjectId:      ProjectId,
		DeleteAtEnd:    DeleteAtEnd,
		CreateAtEnd:    CreateAtEnd,
	}
	command.NewSCW(opts).Init()
}
