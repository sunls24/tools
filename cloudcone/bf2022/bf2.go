package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var waitChange bool
var skip1 bool
var reqInterval = time.Millisecond * 20

var checkCount, vpsCount int

func init() {
	flag.BoolVar(&waitChange, "wc", true, "wait plan change")
	flag.BoolVar(&skip1, "s1", false, "wait plan change")
	flag.IntVar(&checkCount, "cc", 300, "check count")
	flag.IntVar(&vpsCount, "vc", 500, "vps count")
	flag.Parse()
	if skip1 {
		waitChange = false
	}
	log.Println("waitChange:", waitChange)
	log.Println("skip1:", skip1)
	log.Println("checkCount:", checkCount)
	log.Println("vpsCount:", vpsCount)
}

func main() {
	// fmt.Println("start:", time.Now())
	for i := 0; i < checkCount; i++ {
		go func() {
			for {
				checkApi()
				<-time.After(time.Second)
			}
		}()
		<-time.After(reqInterval)
	}
	// fmt.Println("end:", time.Now())
	time.Sleep(time.Hour)
}

func checkApi() {
	if startVps {
		time.Sleep(time.Hour)
		return
	}

	http.DefaultClient.Timeout = 3 * time.Second
	resp, err := http.Get("https://app.cloudcone.com/blackfriday/offers")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var bodyMap = make(map[string]interface{})
	if err = json.Unmarshal(body, &bodyMap); err != nil {
		return
	}

	status, ok := bodyMap["status"].(float64)
	if !ok || status != 1 {
		return
	}
	data, ok := bodyMap["__data"].(map[string]interface{})
	if ok {
		vpsData, ok := data["vps_data"].(map[string]interface{})
		if ok {
			if reflect.DeepEqual(vpsData, last) {
				return
			}
			var oldLast = last
			last = vpsData
			fmt.Println("------------")
			var index int
			for _, v := range vpsData {
				index++
				if skip1 && index == 1 {
					continue
				}
				ins := parseIns(v.(map[string]interface{}))
				ramF := parseRam(ins.ram)
				if ramF < 0.5 {
					continue
				}
				if ins.cpu > 1 && ramF < 1 {
					continue
				}
				if ins.price > 9 {
					continue
				}
				fmt.Printf("ID: %.0f\t\tCPU: %.0f\tRAM: %s\tDISK: %.0f\tPrice: %.2f\t\tbandwidth: %s\n", ins.id, ins.cpu, ins.ram, ins.disk, ins.price, ins.bandwidth)
				if oldLast != nil || !waitChange {
					fmt.Println("------------")
					log.Println("start vps: ", fmt.Sprintf("%+v", ins))
					for i := 1; i < vpsCount; i++ {
						go func() {
							for {
								vps(ins)
								<-time.After(time.Second)
							}
						}()
						<-time.After(reqInterval)
					}
				}
			}
			fmt.Println()
		}
		return
	}
	log.Println("check: not found data body")
}

type ins struct {
	id        float64
	cpu       float64
	ram       string
	disk      float64
	price     float64
	bandwidth string
}

func parseIns(m map[string]interface{}) ins {
	return ins{
		id:        m["id"].(float64),
		cpu:       m["cpu"].(float64),
		ram:       m["ram"].(string),
		disk:      m["disk"].(float64),
		price:     m["usd_price"].(float64),
		bandwidth: m["bandwidth"].(string),
	}
}

var last map[string]interface{}
var startVps bool

func vps(ins ins) {
	const curl = `curl -s 'https://app.cloudcone.com/ajax/vps' \
	-H 'authority: app.cloudcone.com' \
	-H 'accept: application/json, text/javascript, */*; q=0.01' \
	-H 'accept-language: zh-CN,zh;q=0.9' \
	-H 'cache-control: no-cache' \
	-H 'content-type: multipart/form-data; boundary=----WebKitFormBoundaryDtkNU1LsaO0rk65r' \
	-H 'cookie: CCM19=d6nc3tqei51nbv6p2fs2hpiq58; tz=Asia/Shanghai; ref=6079; crisp-client%2Fsession%2Fb4a6582f-f407-4054-b73c-d6e4bf698b1e=session_11e408e7-7343-4c42-ab63-dd29ef2661a4; crisp-client%2Fsession%2Fb4a6582f-f407-4054-b73c-d6e4bf698b1e%2Fdbd21c0fb140422cd05370d5e1c751106aaffce083cdf53dace298fba1eba534=session_11e408e7-7343-4c42-ab63-dd29ef2661a4; bfc_22=1' \
	-H 'origin: https://app.cloudcone.com' \
	-H 'pragma: no-cache' \
	-H 'sec-ch-ua: "Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "macOS"' \
	-H 'sec-fetch-dest: empty' \
	-H 'sec-fetch-mode: cors' \
	-H 'sec-fetch-site: same-origin' \
	-H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36' \
	-H 'x-requested-with: XMLHttpRequest' \
	--data-raw $'------WebKitFormBoundaryDtkNU1LsaO0rk65r\r\nContent-Disposition: form-data; name="os"\r\n\r\n1017\r\n------WebKitFormBoundaryDtkNU1LsaO0rk65r\r\nContent-Disposition: form-data; name="hostname"\r\n\r\narch-cc\r\n------WebKitFormBoundaryDtkNU1LsaO0rk65r\r\nContent-Disposition: form-data; name="plan"\r\n\r\n#####\r\n------WebKitFormBoundaryDtkNU1LsaO0rk65r\r\nContent-Disposition: form-data; name="method"\r\n\r\nprovision\r\n------WebKitFormBoundaryDtkNU1LsaO0rk65r\r\nContent-Disposition: form-data; name="_token"\r\n\r\nTQhMa9jl\r\n------WebKitFormBoundaryDtkNU1LsaO0rk65r--\r\n' \
	--compressed`
	startVps = true
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	cmd := exec.CommandContext(timeout, "bash", "-c", strings.Replace(curl, "#####", strconv.Itoa(int(ins.id)), 1))
	output, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(output))
	if len(out) != 0 {
		log.Println(out)
	}
	if err != nil {
		log.Println(err)
	}
}

func parseRam(ram string) float64 {
	ram = ram[:strings.Index(ram, " ")]
	ret, err := strconv.ParseFloat(ram, 64)
	if err != nil {
		log.Println(err)
	}

	return ret
}
