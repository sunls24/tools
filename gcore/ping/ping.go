package ping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gcore/ip"
	"gcore/utils"
	"io"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

const (
	ipBatch = "http://ip-api.com/batch?fields=countryCode,regionName,query"
)

var (
	pingCount int

	resultList []result
)

type result struct {
	query queryResult
	stat  *ping.Statistics
}

func init() {
	http.DefaultClient.Timeout = 15 * time.Second
}

func PingGcoreIp(count int) {
	pingCount = count
	log.Println("ping count:", pingCount)

	ipList := ip.GetGcoreIp().Addresses

	ipListSp := make([][]string, 1) // 对ip列表进行分组，100为一组

	spIndex := 0
	ipListSp[spIndex] = make([]string, 0, 100)

	for _, ip := range ipList {
		ip = ip[:len(ip)-3]
		ipListSp[spIndex] = append(ipListSp[spIndex], ip)
		if len(ipListSp[spIndex]) == 100 {
			spIndex++
			ipListSp = append(ipListSp, make([]string, 0, 100))
		}
	}

	queryCh := make(chan []queryResult, len(ipListSp))

	var wg sync.WaitGroup
	wg.Add(len(ipListSp))
	for _, sp := range ipListSp {
		go func(sp []string) {
			reqBody, err := json.Marshal(sp)
			utils.Check(err)

			var resp *http.Response
			for {
				resp, err = http.Post(ipBatch, "text/plain;charset=UTF-8", bytes.NewReader(reqBody))
				if err != nil {
					log.Println("retry batch.")
					<-time.After(time.Second)
					continue
				}
				break
			}

			body, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			utils.Check(err)

			var queryList = make([]queryResult, 0, 100)
			err = json.Unmarshal(body, &queryList)
			utils.Check(err)

			queryCh <- queryList
			wg.Done()
		}(sp)
		<-time.After(time.Second)
	}

	var finish = make(chan struct{})
	var resultCh = make(chan result, 100*len(ipListSp))
	go func() {
		wg := sync.WaitGroup{}
		for list := range queryCh {
			for _, q := range list {
				if q.CountryCode != "US" {
					continue
				}
				wg.Add(1)
				go pingIp(q, &wg, resultCh)
			}
		}
		wg.Wait()
		close(finish)
		close(resultCh)
		fmt.Println()
	}()

	go func() {
		for v := range resultCh {
			if v.stat.AvgRtt == 0 {
				continue
			}
			resultList = append(resultList, v)
		}
	}()

	wg.Wait()
	close(queryCh)

	<-finish

	sort.Slice(resultList, func(i, j int) bool {
		return resultList[i].stat.AvgRtt < resultList[j].stat.AvgRtt
	})

	fmt.Println("-------------------------------------------------------------------------------")
	fmt.Println("ip \t\t loss \t avg \t\t min \t\t max \t\t devrtt")

	for _, v := range resultList {
		fmt.Println(v.stat.IPAddr.IP, "\t", fmt.Sprintf("%.0f%%", v.stat.PacketLoss), "\t", v.stat.AvgRtt, "\t", v.stat.MinRtt, "\t", v.stat.MaxRtt, "\t", v.stat.StdDevRtt)
	}
}

func pingIp(query queryResult, wg *sync.WaitGroup, ch chan<- result) {
	defer wg.Done()

	fmt.Print(query.Query, " ")

	pinger, err := ping.NewPinger(query.Query)
	utils.Check(err)

	pinger.Timeout = time.Second * time.Duration(pingCount)
	pinger.Count = pingCount
	if err := pinger.Run(); err != nil {
		return
	}

	ch <- result{
		query: query,
		stat:  pinger.Statistics(),
	}
}

type queryResult struct {
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	Query       string `json:"query"`
}
