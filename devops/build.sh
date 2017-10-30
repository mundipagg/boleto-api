#!/bin/bash
echo "O usuario que esta executando o script: $USER"
#export GOPATH="/home/mundipagg/go"
PROJECTPATH="$GOPATH/src/github.com/mundipagg/boleto-api"
echo "Changing directory to path $PROJECTPATH"
ls -la $PROJECTPATH
cd $PROJECTPATH

echo "Installing dependencies with glide"

echo "Starting build"
go build  -o ./devops/boleto-api -v
