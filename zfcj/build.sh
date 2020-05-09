#!/usr/bin/env zsh

rm -f cmd/bindata/*
go-bindata -split -pkg bindata -o ./cmd/bindata ui/dist/...
#go run main.go serve
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o bin/zfcj main.go
docker build -t bin/zfcj .
docker tag zfcj registry.cn-hangzhou.aliyuncs.com/golden/zfcj
docker push registry.cn-hangzhou.aliyuncs.com/golden/zfcj
