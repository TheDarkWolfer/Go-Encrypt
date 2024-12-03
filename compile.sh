#!/bin/bash

# The chosen file
GO_FILE="./encryption.go"
GO_FILE=$1

# The version string
VERSION="42"
VERSION=$2

# Check if the user correctly provided the file to compile
[ "$VERSION" == "42" ] && echo "Not specifying version number"

echo "Building $GO_FILE for Linux"
#Build for Linux and check if it compiled right, and notice the user of the result
go build -ldflags "-X main.Version=$VERSION -X main.BuildDate=$(date +'%Y-%m-%d')" -o ${GO_FILE}_linux ${GO_FILE} 
[ $? -eq 0 ] && echo "Successfully build for Linux !" || echo "Linux build failed !"

echo "Building $GO_FILE for Windows"
#Build for Windows and check if it compiled right, and notice the user of the result
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION -X main.BuildDate=$(date +'%Y-%m-%d')" -o ${GO_FILE}_windows.exe ${GO_FILE}
[ $? -eq 0 ] && echo "Successfully build for Windows !" || echo "Windows build failed !"