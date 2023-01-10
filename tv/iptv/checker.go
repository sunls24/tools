package iptv

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultFirst = `#EXTM3U x-tvg-url="https://sunls.me/d/tv/online.xml"`
)

func init() {
	http.DefaultClient.Timeout = time.Minute
}

func Checker(opts Options) {
	if len(opts.Available) == 0 {
		// iptv-checker
		fmt.Println("start exec iptv-checker...")
		cmd := exec.Command("bash", "-c", fmt.Sprintf(`iptv-checker -k -p %d -o %s %s`, opts.Parallel, opts.Workdir, opts.Source))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println(cmd.Args)
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		opts.Available = fmt.Sprintf("%s/online.m3u", opts.Workdir)
		fmt.Println("end exec iptv-checker")
	}
	availableList := readPlaylist(opts.Available)
	streamMap := make(map[string]int, len(availableList.streamList))
	for i, a := range availableList.streamList {
		url := strings.TrimSpace(a.url)
		if _, ok := streamMap[url]; ok {
			fmt.Println("duplicate:", url)
			continue
		}
		streamMap[url] = i
	}

	sourceList := readPlaylist(opts.Source)

	newList := newPlaylist(defaultFirst, make([]stream, 0, len(availableList.streamList)))
	for _, s := range sourceList.streamList {
		url := strings.TrimSpace(s.url)
		if _, ok := streamMap[url]; ok {
			newList.streamList = append(newList.streamList, s)
		}
	}

	fmt.Println("source:", len(sourceList.streamList))
	fmt.Println("available:", len(availableList.streamList), "|", "stream map:", len(streamMap))
	fmt.Println("new:", len(newList.streamList))

	output, err := os.OpenFile(opts.Output, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	newList.fmt(bufio.NewWriter(output))
	fmt.Println("save m3u:", opts.Output)

	if len(opts.EpgXml) != 0 {
		checkEpgXML(newList.streamList, opts)
	}
}

func checkEpgXML(sList []stream, opts Options) {
	fmt.Println("--------")
	const tvgFlag = `tvg-id="`
	var sMap = make(map[string]struct{}, len(sList))
	for _, s := range sList {
		for _, v := range s.comment {
			if i := strings.Index(v, tvgFlag); i >= 0 {
				v = v[i+len(tvgFlag):]
				if i = strings.Index(v, `"`); i > 0 {
					sMap[v[:i]] = struct{}{}
				}
			}
		}
	}
	fmt.Println("tvg:", len(sMap))

	data := readSource(opts.EpgXml)
	var result []string
	var count int
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}
		if len(result) == 0 && strings.HasPrefix(txt, "<?xml") {
			result = append(result, txt)
		}

		const (
			idFlag      = `id="`
			channelFlag = `channel="`
		)
		if i := strings.Index(txt, idFlag); i > 0 {
			id := txt[i+len(idFlag):]
			id = id[:strings.Index(id, `"`)]
			if _, ok := sMap[id]; ok {
				result = append(result, txt)
				count++
			}
		}

		if i := strings.Index(txt, channelFlag); i > 0 {
			channel := txt[i+len(channelFlag):]
			channel = channel[:strings.Index(channel, `"`)]
			if _, ok := sMap[channel]; ok {
				result = append(result, txt)
			}
		}
	}
	result = append(result, "</tv>")

	output, err := os.OpenFile(opts.EpgOutput, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	cf := bufio.NewWriter(output)
	for _, v := range result {
		cf.WriteString(v)
		cf.WriteByte('\n')
	}
	cf.Flush()
	fmt.Println("match:", count)
	fmt.Println("save epg:", opts.EpgOutput)
}

func readSource(s string) io.ReadCloser {
	var rc io.ReadCloser
	var err error
	if strings.HasPrefix(s, "http") {
		rc, err = download(s)
	} else {
		rc, err = os.OpenFile(s, os.O_CREATE|os.O_RDWR, 0o644)
	}
	if err != nil {
		panic(err)
	}
	return rc
}

func readPlaylist(s string) playlist {
	return parsePlaylist(readSource(s))
}

func parsePlaylist(readCloser io.ReadCloser) playlist {
	defer readCloser.Close()
	var streamList []stream
	var cur *stream

	var first string
	scanner := bufio.NewScanner(readCloser)
	for scanner.Scan() {
		txt := scanner.Text()
		if len(txt) == 0 {
			continue
		}
		if len(first) == 0 && strings.HasPrefix(txt, "#EXTM3U") {
			first = txt
			continue
		}
		if cur == nil {
			cur = newStream()
		}
		if strings.HasPrefix(txt, "#") {
			cur.comment = append(cur.comment, txt)
		} else if strings.HasPrefix(txt, "http") {
			cur.url = txt
			streamList = append(streamList, *cur)
			cur = nil
		}
	}
	return newPlaylist(first, streamList)
}

func download(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	return resp.Body, err
}
