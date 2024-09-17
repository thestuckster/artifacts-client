#!/bin/bash

go mod tidy
rm artifacts-client
go build . 
./artifacts-client
