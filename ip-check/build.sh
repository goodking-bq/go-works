#!/usr/bin/env zsh



export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o bin/ipcheck main.go



export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o bin/ipcheck.exe main.go
