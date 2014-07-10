package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	port = ":8080"
)

var blocked = false

func main() {

	http.HandleFunc("/", httpHandler)

	fmt.Println("Running on port", port)
	panic(http.ListenAndServe(port, nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {

	if blocked {

		http.Error(w, http.StatusText(503), 503)

	} else {

		blocked = true

		cmd := exec.Command("sh", "/home/pi/camera/picture")
		go cmd.Output()

		fileCb := make(chan *os.File)
		timeoutCb := make(chan bool)

		go monitorFile("/home/pi/camera/pic.jpg", fileCb, timeoutCb)

		select {

		case file := <-fileCb:

			buf := make([]byte, 1024*1024*10)

			start := 0

			//Write data to response buffer
			for {

				n, err := file.Read(buf)

				if err != nil && err != io.EOF {

					panic(err)
				}
				if n == 0 {
					break
				}

				_, err = w.Write(buf[start:n])

				if err != nil {

					panic(err)
				}

				start = n
			}

			w.Header().Set("Content-Type", "image/jpg")
			blocked = false

		case <-timeoutCb:

			blocked = false
			http.Error(w, http.StatusText(500), 500)
		}
	}
}

func monitorFile(file string, cb chan *os.File, timeout chan bool) {

	ttl := 50

	for {

		if ttl == 0 {

			timeout <- true
		}

		time.Sleep(100 * time.Millisecond)
		file, err := os.Open("/home/pi/camera/pic.jpg")

		if err == nil && file != nil {

			defer file.Close()

			cb <- file
		}

		ttl -= 1
	}
}
