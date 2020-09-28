#!/usr/bin/env bash

protoc -I=./ --gofast_out=./ ./barrage/barrage.proto
