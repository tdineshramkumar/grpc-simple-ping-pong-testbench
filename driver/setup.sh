#!/bin/bash
echo "Building c++ binaries"
make -f ../cpp/Makefile -B -C ../cpp all
echo "Copying them into current directory"
cp ../cpp/client cpp-client
cp ../cpp/server cpp-server
echo "Building go binaries"
go build ../pingpong/client/*.go
go build ../pingpong/server/*.go
echo "Copying them into current directory"
cp ../pingpong/server/server go-server
cp ../pingpong/client/client go-client

