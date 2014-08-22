package main

import (
	"flag"
	"time"
)

var (
	eventCode      = flag.String("event", "s9pgpX", "Event code for server")
	directoryCode  = flag.String("directory", "HT6HWWY7KWVNHCWMTWXMYXRVFXUTJXZZ", "Directory code for amazon")
	zulipStream    = flag.String("zulipStream", "test-stream", "Stream for messages")
	zulipSubject   = flag.String("zulipSubject", "photobooth", "Subject for messages")
	listenInterval = flag.Int("interval", 60, "time interval to listen on server")
	listener       = &PhotoboothListener{}
	imager         = &AmazonImageHandler{}
	zulip          = &ZulipBot{}
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

	i := time.Duration(*listenInterval)
	listener.Interval = time.Second * i
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
