#!/usr/bin/env zsh

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o bin/ipcheck main.go
scp -P 59878 -i ~/.ssh/platform bin/ipcheck platform@ipcheck.2xi.com:/tmp
ssh -p 59878 -i ~/.ssh/platform platform@ipcheck.2xi.com "sudo systemctl stop ipcheck && cp -f /tmp/ipcheck /data/ipcheck/ &&sudo  systemctl start ipcheck"

export CGO_ENABLED=0
export GOOS=windows
export GOARCH=amd64
go build -o bin/ipcheck.exe main.go
