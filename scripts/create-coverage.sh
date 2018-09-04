#!/bin/bash

. ./set-env.sh
go test -coverprofile=c.out
sed -i "s/_$(pwd|sed 's/\//\\\//g')/./g" c.out
go tool cover -html=c.out -o=coverage.html
rm c.out
echo "Generated coverage.html!"