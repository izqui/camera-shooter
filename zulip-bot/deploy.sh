#!/bin/bash

GOARCH=arm GOOS=linux go build
scp zulip-bot pi@10.0.5.194:zulip
