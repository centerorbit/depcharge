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
https://github.com/golang/go/issues/22430#issuecomment-414668599

```
go test -coverprofile=c.out
sed "s/_$(pwd|sed 's/\//\\\//g')/./g" c.out > c.out
go tool cover -html=c.out -o=c.html
```

## Semantic Versioning
https://semver.org/

Given a version number MAJOR.MINOR.PATCH, increment the:

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

Examples:

* Regular release: v1.0.0
* [Pre-release versions](https://semver.org/#spec-item-9): v1.0.0-rc.2
* [Build metadata](https://semver.org/#spec-item-10): v1.0.0-rc.1+b78be23