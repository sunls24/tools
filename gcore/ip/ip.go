package ip

import (
	"encoding/json"
	"fmt"
	"gcore/utils"
	"io"
	"net/http"
	"strings"
)

const (
	ipApi = "https://api.gcorelabs.com/cdn/public-ip-list"
)

func GetGcoreIp() *IpList {
	resp, err := http.Get(ipApi)
	utils.Check(err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	utils.Check(err)

	var ipList IpList
	utils.Check(json.Unmarshal(body, &ipList))
	return &ipList
}

type IpList struct {
	Addresses   []string `json:"addresses"`
	AddressesV6 []string `json:"addresses_v6"`
}

func PrintIp(list []string) {
	fmt.Println(strings.Join(list, "\n"))
}
