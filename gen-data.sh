#!/bin/bash

go get -u github.com/jteeuwen/go-bindata/...;
go-bindata -ignore="data.go" -o data/assets.go -pkg data data/...;
