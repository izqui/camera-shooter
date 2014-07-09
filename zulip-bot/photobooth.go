package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	PHOTOBOOTH_BASEURL = "http://yetibellyblue.com"
)

var (
	httpClient = &http.Client{}
)

type LastPictureCallback chan int
type Listener struct {
	Interval  time.Duration
	LastValue int
	callback  LastPictureCallback
	event     string
	listening bool
}

func (l *Listener) ListenForChangesInEvent(event string, cb LastPictureCallback) {

	l.callback = cb
	l.listening = true
	l.event = event

	if l.Interval == 0 {
		l.Interval = 15 * time.Second
	}

	l.LastValue = -1

	fmt.Println("Listening for Photobooth:", l.event, "every", l.Interval)
	go l.listen()
}

func (l *Listener) StopListening() {

	l.listening = false
}

func (l *Listener) listen() {

	for l.listening {

		i := l.getLastPicture()

		if i > -1 {

			if l.LastValue < 0 {

				fmt.Println("Initial value: ", i)
				l.LastValue = i
			} else {

				if i > l.LastValue {

					l.LastValue = i
					l.callback <- i
				}
			}
		}

		time.Sleep(l.Interval)
	}
}

func (l *Listener) getLastPicture() int {

	url := fmt.Sprintf("%s/getEventData?eventcode=%s", PHOTOBOOTH_BASEURL, l.event)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Cookie", COOKIE)
	res, err := httpClient.Do(req)

	if err == nil {

		data, err := ioutil.ReadAll(res.Body)
		if err == nil {

			d := make(map[string]interface{})
			json.Unmarshal(data, &d)

			if d["EventData"] != nil {

				eventData := d["EventData"].(map[string]interface{})
				if eventData["last"] != nil {

					last := eventData["last"].(string)

					i, err := strconv.Atoi(last)
					if err == nil {

						return i
					}
				}
			}
		}
	}

	return -1
}
