# Contributing
* [GitLab is where development & pipelines occur](https://gitlab.com/centerorbit/depcharge)
* [GitHub is the public-access repository](https://github.com/centerorbit/depcharge)

## Code Setup
1. Clone repo
1. Done! (DepCharge uses Go Modules, and so all dependencies will be installed automatically.)

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
sed -i "s/_$(pwd|sed 's/\//\\\//g')/./g" c.out
go tool cover -html=c.out -o=c.html
```

### Strict cover
The CI pipeline will ensure that at _least_ 80% code coverage exists. The following env will ensure this strict coverage percentage is checked. The `go test -cover` command will fail if your coverage is below 80% and this env is set:

`export COVER_STRICT=true`

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