#!/bin/bash

GOARCH=arm GOOS=linux go build
scp camera pi@10.0.5.195:camera/shooter
