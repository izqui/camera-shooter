package main

import (
	"flag"
	"time"
)

var (
	eventCode     = flag.String("event", "nG8z3q", "Event code for server")
	directoryCode = flag.String("directory", "IIDSJ7YNOO32UPW5T5QWFDX33CIMTQYY", "Directory code for amazon")
	zulipStream   = flag.String("zulipStream", "test-stream", "Stream for messages")
	zulipSubject  = flag.String("zulipSubject", "photobooth", "Subject for messages")
	listener      = &PhotoboothListener{}
	imager        = &AmazonImageHandler{}
	zulip         = &ZulipBot{}
)

func init() {

	flag.Parse()

	zulip.Stream = *zulipStream
	zulip.Subject = *zulipSubject

	listener.Event = *eventCode

	imager.Directory = *directoryCode
}

func main() {

	var listenerCallback = make(LastPictureCallback)
	var imagerCallback = make(MarkdownImageTextCallback)

	listener.Interval = 3 * time.Second
	listener.ListenForChanges(listenerCallback)

	imager.Callback = imagerCallback

	//imager.GetImageMarkdownRepresentation(CallbackData{Start: 40, End: 55})

	for {

		select {

		case d := <-listenerCallback:
			imager.GetImageMarkdownRepresentation(d)

		case m := <-imagerCallback:
			zulip.SendMessage(m)
		}
	}
}
