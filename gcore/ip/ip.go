package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	ipApi = "https://api.gcorelabs.com/cdn/public-ip-list"
)

func GetGcoreIp() *IpList {
	resp, err := http.Get(ipApi)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	var ipList IpList
	if err = json.Unmarshal(body, &ipList); err != nil {
		log.Panicln(err)
	}
	return &ipList
}

type IpList struct {
	Addresses   []string `json:"addresses"`
	AddressesV6 []string `json:"addresses_v6"`
}

func PrintIp(list []string) {
	fmt.Println(strings.Join(list, "\n"))
}
