package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
)

var (
	secret string
	target string
)

func init() {
	flag.StringVar(&secret, "e", "", "secret key")
	flag.StringVar(&target, "t", "", "hmac sha1 target string")
	flag.Parse()
}

func main() {
	if len(secret) == 0 {
		fmt.Println("secret key is empty")
		flag.PrintDefaults()
		return
	}
	fmt.Println(hmacSha12(secret, target))
}

func hmacSha12(secret, payload string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedMAC := mac.Sum(nil)
	return hex.EncodeToString(expectedMAC)
}
