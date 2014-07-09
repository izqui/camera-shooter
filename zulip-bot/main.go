package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	eventCode     = flag.String("event", "nG8z3q", "Event code for server")
	directoryCode = flag.String("directory", "IIDSJ7YNOO32UPW5T5QWFDX33CIMTQYY", "Directory code for amazon")
	zulipStream   = flag.String("zulipStream", "test-stream", "Stream for messages")
	zulipSubject  = flag.String("zulipSubject", "photobooth", "Subject for messages")
	listener      = &Listener{}
	zulip         = &ZulipBot{}
)

func init() {

	flag.Parse()
	zulip.Stream = *zulipStream
	zulip.Subject = *zulipSubject
}

func main() {

	var callback = make(LastPictureCallback)

	listener.Interval = 10 * time.Second
	listener.ListenForChangesInEvent(*eventCode, callback)

	fmt.Println(zulip.SendMessage("the Test"))

	for {

		select {

		case i := <-callback:

			fmt.Println("Changed", i)
		}
	}
}
