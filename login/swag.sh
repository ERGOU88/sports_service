#!/bin/bash

swag init;
sed 's/"type": "file"/"type": "string"/g' ./docs/docs.go > ./docs/docs.go.tmp;
sed 's/multipart\/form-data/application\/x-www-form-urlencoded/g' ./docs/docs.go.tmp > ./docs/docs.go;
go run main.go
