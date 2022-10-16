#!/bin/bash

rm Robintris-Linux-x64
echo "Delete of Robintris-Linux-x64 finished. Building new one."
export GOARCH='amd64'
go build -o Robintris-Linux-x64 -tags=ebitenginesinglethread -ldflags="-s -w" main.go

rm Robintris-Linux-x86
echo "Delete of Robintris-Linux-x86 finished. Building new one."
export GOARCH='386'
go build -o Robintris-Linux-x86 -tags=ebitenginesinglethread -ldflags="-s -w" main.go

rm Robintris-Linux-arm64
echo "Delete of Robintris-Linux-arm64 finished. Building new one."
export GOARCH='arm64'
go build -o Robintris-Linux-arm64 -tags=ebitenginesinglethread -ldflags="-s -w" main.go

rm Robintris-Windows-x64.exe
echo "Delete of Robintris-Windows-x64.exe finished. Building new one."
export GOOS='windows'
export GOARCH='amd64'
go build -o Robintris-Windows-x64.exe -tags=ebitenginesinglethread main.go

rm Robintris-Windows-x86.exe
echo "Delete of Robintris-Windows-x86.exe finished. Building new one."
export GOARCH='386'
go build -o Robintris-Windows-x86.exe -tags=ebitenginesinglethread main.go

rm Robintris-Windows-arm64.exe
echo "Delete of Robintris-Windows-arm64.exe finished. Building new one."
export GOARCH='arm64'
go build -o Robintris-Windows-arm64.exe -tags=ebitenginesinglethread main.go

rm Robintris.wasm
echo "Delete of Robintris.wasm finished. Building new one."
export GOOS='js'
export GOARCH='wasm'
go build -o Robintris.wasm -tags=ebitenginesinglethread -ldflags="-s -w" main.go
