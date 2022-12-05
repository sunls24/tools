package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var interval string

func init() {
	flag.StringVar(&interval, "t", "1m", "check interval")
	flag.Parse()
}

func main() {
	if interval == "" {
		flag.PrintDefaults()
		return
	}
	wait, err := time.ParseDuration(interval)
	if err != nil {
		log.Panicln(err)
	}
	sendTgMsg("test send msg by first run")
	for {
		check()
		<-time.After(wait)
	}
}

func check() {
	resp, err := http.Get("https://app.cloudcone.com/blackfriday/offers")
	if err != nil {
		log.Println("http.Get:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body:", err)
		return
	}
	var bodyMap = make(map[string]interface{})
	if err = json.Unmarshal(body, &bodyMap); err != nil {
		log.Println("unmarshal:", err)
		return
	}

	status, ok := bodyMap["status"].(float64)
	if !ok || status != 0 {
		log.Println("status code != 0:", fmt.Sprintf("%+v", bodyMap))
		return
	}
	data, ok := bodyMap["__data"].(map[string]interface{})
	if ok {
		log.Println("find data: ", fmt.Sprintf("%+v", data))
		vpsData, ok := data["vps_data"].(map[string]interface{})
		if ok {
			for k := range vpsData {
				ki, err := strconv.Atoi(k)
				if err != nil && ki != 0 {
					sendTgMsg(fmt.Sprintf("find vps data body: %+v", vpsData))
				}
			}
		}
		return
	}
	log.Println("check: not found data body")
}

func sendTgMsg(msg string) {
	curl := fmt.Sprintf(`curl -k --data chat_id="1969557829" --data "text=%s" "https://api.telegram.org/bot5984183633:AAGYRJlP12lqyc8rMlnKKdizguZgrKqiqPc/sendMessage"`, msg)
	cmd := exec.Command("bash", "-c", curl)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	out := strings.TrimSpace(string(output))
	log.Println("> sendTgMsg:", out)
}
