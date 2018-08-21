# Contributing

## Formatting code

## Test Coverage

go test -coverprofile=c.out
(edit out to be `./` relative path) 
go tool cover -html=c.out -o=result.html