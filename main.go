package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	port = ":8080"
)

func main() {

	http.HandleFunc("/", httpHandler)

	fmt.Println("Running on port", port)
	panic(http.ListenAndServe(port, nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("request")

	cmd := exec.Command("sh", "/home/pi/camera/picture.sh")
	_, err := cmd.Output()

	if err == nil {

		file, err := os.Open("/home/pi/camera/pic.jpg")

		if err == nil {

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

		} else {

			panic(err)
		}
		defer func() {

			if err := file.Close(); err != nil {

				panic(err)
			}
		}()

	}
}
