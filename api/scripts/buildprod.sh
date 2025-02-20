#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webjson-public-api ./cmd/public
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webjson-web-api ./cmd/web

