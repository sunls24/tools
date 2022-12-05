package command

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

type SCW struct {
	insUUID string
	opts    Options
}

type Options struct {
	StartCount     int           // 开启次数
	StartInterval  time.Duration // 开启间隔
	StatusInterval time.Duration // 检测状态间隔
	PrintCommand   bool          // 打印命令
	ProjectId      string

	DeleteAtEnd bool // 当结束时删除实例
	CreateAtEnd bool // 当结束时删除并重建实例，需要 DeleteAtEnd 为 True
}

func NewSCW(opts Options) *SCW {
	return &SCW{opts: opts}
}

func (scw *SCW) Init() {
	scw.SetInsUUID()
	if len(scw.insUUID) == 0 {
		scw.CreateINS()
	}

	if len(scw.insUUID) == 0 {
		log.Panicln("init: ins uuid is empty")
	}

	scw.schedule()
}

func (scw *SCW) GetStatus() string {
	var cmd = "scw instance server list | sed -n '2p' | awk '{print $4}'"
	output, err := scw.execCommand(cmd)
	if err != nil {
		log.Panicln("GetStatus")
	}
	return strings.TrimSpace(output)
}

func (scw *SCW) SetInsUUID() {
	const cmd = "scw instance server list | sed -n '2p' | awk '{print $1}'"
	output, err := scw.execCommand(cmd)
	if err != nil {
		log.Panicln("SetInsUUID: ", err)
	}
	if len(output) == 0 {
		log.Println("SetInsUUID output is empty")
		return
	}
	scw.insUUID = output
	log.Println("SetInsUUID:", scw.insUUID)
}

func (scw *SCW) CreateINS() {
	var cmd = fmt.Sprintf("scw instance server create type=STARDUST1-S zone=fr-par-1 image=arch_linux root-volume=l:10G name=scw-arch-par1 ip=none ipv6=true project-id=%s", scw.opts.ProjectId)
	_, err := scw.execCommand(cmd)
	if err == nil {
		scw.SetInsUUID()
	}
}

func (scw *SCW) StartINS() {
	if len(scw.insUUID) == 0 {
		log.Println("StartINS: ins uuid is empty")
		return
	}
	var cmd = fmt.Sprintf("scw instance server start %s", scw.insUUID)
	_, _ = scw.execCommand(cmd)
}

func (scw *SCW) DeleteINS() {
	if len(scw.insUUID) == 0 {
		log.Println("StartINS: ins uuid is empty")
		return
	}
	var cmd = fmt.Sprintf("scw instance server delete %s", scw.insUUID)
	_, _ = scw.execCommand(cmd)
}

func (scw *SCW) execCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(output))
	scw.errorCheck(command, out, err)
	return out, err
}

func (scw *SCW) errorCheck(cmd, output string, err error) {
	if !scw.opts.PrintCommand && err == nil {
		return
	}
	log.Println("> command:", cmd)
	log.Println("> output:", output)
	if err == nil {
		return
	}
	log.Println("⬆⬆⬆ command exec error:", err)
}

func (scw *SCW) schedule() {
	log.Println("start schedule:", scw.insUUID)
	log.Println("now status =", scw.GetStatus())
	for i := 1; i <= scw.opts.StartCount; {
		scw.waitStatus(Archived)

		scw.StartINS()
		log.Println("status is archived, now start ins count", i)
		i++
		<-time.After(scw.opts.StartInterval)
	}

	if scw.opts.DeleteAtEnd {
		log.Println("start count end, prepare delete ins")
		scw.waitStatus(Archived)
		scw.DeleteINS()
		if scw.opts.CreateAtEnd {
			log.Println("delete ins success, wait", scw.opts.StartInterval)
			<-time.After(scw.opts.StartInterval)
			log.Println("create new ins")
			scw.CreateINS()
			scw.schedule()
		}
	}
}

const (
	Empty    = ""
	Starting = "starting"
	Archived = "archived"
)

func (scw *SCW) waitStatus(target string) {
	for scw.GetStatus() != target {
		log.Println("wait status =", target, ", now status =", scw.GetStatus())
		<-time.After(scw.opts.StatusInterval)
	}
}
