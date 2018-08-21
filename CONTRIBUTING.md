# Contributing

## Code Setup
1. Clone repo
1. `go get -t ./...`

(-t includes test packages)

## Formatting code

This command will auto-clean the formatting of project code:

`gofmt -s -w *.go`

## Testing and Coverage

On most days:
`go test -cover`

### Generate coverage report:
`go test -coverprofile=c.out`

Edit c.out to be `./` relative path, at least until I figure out how to properly work with it.
https://github.com/Masterminds/glide/issues/43

`go tool cover -html=c.out -o=result.html`