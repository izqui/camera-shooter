package main

import (
	"fmt"
	"net/http"
)

const (
	S3_BASEURL      = "http://photosaurus.s3.amazonaws.com"
	MARKDOWN_FORMAT = "[Image](%s)"

	THUMB_FORMAT = "%s/%s/%sthumb.jpg"
	BIG_FORMAT   = "%s/%s/%slarge.jpg"
)

type MarkdownImageTextCallback chan string
type AmazonImageHandler struct {
	Directory string
	Callback  MarkdownImageTextCallback
}

func (a *AmazonImageHandler) GetImageMarkdownRepresentation(data CallbackData) {

	for i := data.Start + 1; i <= data.End; i++ {

		go a.checkImage(fmt.Sprintf("%d", i))
	}
}

func (a *AmazonImageHandler) checkImage(i string) {

	url := fmt.Sprintf(THUMB_FORMAT, S3_BASEURL, a.Directory, i)

	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 200 {

		a.Callback <- fmt.Sprintf(MARKDOWN_FORMAT, fmt.Sprintf(BIG_FORMAT, S3_BASEURL, a.Directory, i))
	}
}
