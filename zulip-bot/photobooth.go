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

type CallbackData struct{ Start, End int }

type LastPictureCallback chan CallbackData

type PhotoboothListener struct {
	Interval time.Duration

	Event string

	LastValue int
	callback  LastPictureCallback
	listening bool
}

func (l *PhotoboothListener) ListenForChanges(cb LastPictureCallback) {

	l.callback = cb
	l.listening = true

	if l.Interval == 0 {
		l.Interval = 15 * time.Second
	}

	l.LastValue = -1

	fmt.Println("Listening for Photobooth:", l.Event, "every", l.Interval)
	go l.listen()
}

func (l *PhotoboothListener) StopListening() {

	l.listening = false
}

func (l *PhotoboothListener) listen() {

	for l.listening {

		i := l.getLastPicture()

		fmt.Println(i)

		if i > -1 {

			if l.LastValue < 0 {

				fmt.Println("Initial value: ", i)
				l.LastValue = i
			} else {

				if i > l.LastValue {

					l.callback <- CallbackData{Start: l.LastValue, End: i}
					l.LastValue = i

				}
			}
		}

		time.Sleep(l.Interval)
	}
}

func (l *PhotoboothListener) getLastPicture() int {

	url := fmt.Sprintf("%s/getEventData?eventcode=%s", PHOTOBOOTH_BASEURL, l.Event)
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
