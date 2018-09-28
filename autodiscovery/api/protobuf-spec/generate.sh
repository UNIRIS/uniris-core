#!/bin/sh

protoc --go_out=plugins=grpc:. discovery.proto