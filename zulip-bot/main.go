package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	eventCode     = flag.String("event", "nG8z3q", "Event code for server")
	directoryCode = flag.String("directory", "IIDSJ7YNOO32UPW5T5QWFDX33CIMTQYY", "Directory code for amazon")
	listener      = &Listener{}
)

func init() {

	flag.Parse()
}

func main() {

	var callback = make(LastPictureCallback)

	listener.Interval = 10 * time.Second
	listener.ListenForChangesInEvent(*eventCode, callback)

	for {

		select {

		case i := <-callback:

			fmt.Println("Changed", i)
		}
	}
}
