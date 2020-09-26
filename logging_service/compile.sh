#!/bin/sh

mkdir -p /go/src/github.com/Jakatut/NAD-A3 \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

export GOPATH=/go
export PATH=$GOPATH/bin:$PATH
export APP=$GOPATH/src/$REPO
export PORT_MAD=8080

cp -r ./app/vendor/** $APP
cp -r  ./app/*.go $APP
cp -r  ./app/messages $APP
cp -r ./app/routes $APP
cp -r ./app/handlers $APP

cd $APP
go build -o NAD-A3
