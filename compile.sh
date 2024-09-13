#!/bin/bash

# The chosen file
GO_FILE=$1

# The version string
VERSION=$2

# Check if the user correctly provided the file to compile
[ -z "$GO_FILE" ] && echo "First argument should be the go file to compile" && exit 1
[ -z "$VERSION" ] && echo "Second argument should be the version number" && exit 1

echo "Building $GO_FILE for Linux"
#Build for Linux and check if it compiled right, and notice the user of the result
go build -ldflags "-X main.Version=$VERSION -X main.BuildDate=$(date +'%Y-%m-%d')" -o ${GO_FILE}_linux ${GO_FILE} 
[ $? -eq 0 ] && echo "Successfully build for Linux !" || echo "Linux build failed !"

echo "Building $GO_FILE for Windows"
#Build for Windows and check if it compiled right, and notice the user of the result
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION -X main.BuildDate=$(date +'%Y-%m-%d')" -o ${GO_FILE}_windows.exe ${GO_FILE}
[ $? -eq 0 ] && echo "Successfully build for Windows !" || echo "Windows build failed !"