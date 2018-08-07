# DepCharge
[![pipeline status](https://gitlab.com/centerorbit/depcharge/badges/master/pipeline.svg)](https://gitlab.com/centerorbit/depcharge/commits/master) [![coverage report](https://gitlab.com/centerorbit/depcharge/badges/master/coverage.svg)](https://gitlab.com/centerorbit/depcharge/commits/master)

DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once.

## Introduction
A medium-to-large sized project (especially when using a microservice architecture) will consist of 3 or more separate repositories, and rely on a variety of package managers depending on the various languages chosen for each service. Typically, these repos must be managed, tracked, and released in some semblance of unison so that the dependant service calls can be understood and responded to appropriately.

For small (to single) teams, a single developer will often need to propagating and perform the same commands across all relevant services. This is a tedious, manual, and error-prone process that can occur every release.

DepCharge is designed to help fix that.

By creating a YAML file that describes all of your project(s) dependencies, you can then execute commands across all of them simultaneously.

## How it Works & Usage
All of the examples here are just that: examples. DepCharge is designed to be as flexible as possible, so if you happen to use tools other than what's listed, they should work as well!

DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once.

Usage: `depcharge --kind=<kind> [--labels=<comma-separated,inherited>] [OPTIONS...] COMMAND [ARGS...]`

### Features:
* Supports arbitrary params, whatever 'params: key: value' pairs you want
* Built-in mustache templating, allows you to parametrize your commands
* Supports YAML anchors

### Description:
`depcharge` will read the `dep.yml` file in the current working directory, and
perform all commands relative to that location.

### Example `dep.yml`:
```
deps:
    - name: frontend
      kind: git
      location: ./app/frontend
      labels:
        - public
      params:
        repo: git@example.com:frontend.git
      deps:
        - name: vue.js
          kind: npm
    - name: backend
      kind: git
      location: ./app/backend
      labels:
        - api
      params:
        repo: git@example.com:backend.git
      deps:
        - name: lumen
          kind: composer
```

### Primary Commands:

 --kind		Is the top-level filter that's applied, opperations are run based on 'kind'

 --labels	Comma separated list of labels to filter by, inherited from parents

### Available Options:

 --help			Shows this message

 --dryrun		Prints out intended command without executing it

 --exclusive	(default) For a match to be found, it must contain at least all provided labels

 --inclusive   	For a match to be found, it must contain at least one of the provided labels

### Example commands:

Will run `git clone <location>` across all git dependencies:

	depcharge --kind=git clone {{repo}} {{location}}
	
Will run `git status` across all git dependencies:

	depcharge --kind=git status
	
Will run `npm install` across any npm dependencies that have the label 'public':

	depcharge --kind=npm --labels=public install
	
Will run `composer install` across any composer dependencies that have either the label 'api', or 'soap':

	depcharge --kind=composer --inclusive --labels=api,soap install
	
And much more!


## Additional Resources
* https://mustache.github.io/

## Contributing
* [GitLab is where development & pipelines occur](https://gitlab.com/centerorbit/depcharge)
* [GitHub is the public-access repository](https://github.com/centerorbit/depcharge)


## License

- [LICENSE](LICENSE) (Expat/[MIT License][MIT])

[MIT]: http://www.opensource.org/licenses/MIT "The MIT License (MIT)"