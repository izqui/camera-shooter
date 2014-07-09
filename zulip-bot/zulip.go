package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

var (
	zulipHttpClient = &http.Client{}
)

const (
	ZULIP_BASEURL = "https://api.zulip.com/v1"
)

type ZulipBot struct {
	Stream  string
	Subject string
}

func (z *ZulipBot) SendMessage(message string) error {

	data := url.Values{}
	data.Set("type", "stream")
	data.Add("to", z.Stream)
	data.Add("subject", z.Subject)
	data.Add("content", message)

	encoded := data.Encode()

	urlstr := fmt.Sprintf("%s/messages?%s", ZULIP_BASEURL, encoded)

	req, _ := http.NewRequest("POST", urlstr, nil)
	req.Header.Add("Authorization", authHeader())
	res, err := zulipHttpClient.Do(req)

	return err

}

func authHeader() string {

	baseString := fmt.Sprintf("%s:%s", ZULIP_LOGIN, ZULIP_PASSWD)
	encoded := base64.StdEncoding.EncodeToString([]byte(baseString))

	return fmt.Sprintf("Basic %s", encoded)
}
