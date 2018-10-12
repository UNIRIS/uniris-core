#!/bin/sh

protoc --go_out=plugins=grpc:. wallet.proto