# DepCharge

DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once.

* [GitHub is the public-access repository](https://github.com/centerorbit/depcharge)
* [GitLab is where development & pipelines occur](https://gitlab.com/centerorbit/depcharge)

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

 --dryrun		*NOT YET IMPLEMENTED!*

 --exclusive	(default) For a match to be found, it must contain at least all provided labels

 --inclusive   	For a match to be found, it must contain at least one of the provided labels

### Example commands:

Will run `git clone <location>` across all git dependencies:

	depcharge --kind=git clone {{location}}
	
Will run `git status` across all git dependencies:

	depcharge --kind=git status
	
Will run `npm install` across any npm dependencies that have the label 'public':

	depcharge --kind=npm --labels=public install
	
Will run `composer install` across any composer dependencies that have either the label 'api', or 'soap':

	depcharge --kind=composer --inclusive --labels=api,soap install
	
And much more!


## ToDo:
### Before v1.0
* Setup sample project for example and test
* Build Pipeline
* Dryrun flag (I think we need a struct for actionHandlers)

### Wish List
1. Test support of YAML anchors
2. Secrets managment?
3. docker-compose handler
4. Implement channels / goroutines
5. Mechanism to perform _other_ types of commands to 'kind's rather than `kind`
6. Need a "literal" flag, to turn off templating?
7. Find a way to "stream" output to terminal?
8. Test Coverage
9. Better Logging & Verbose flag


## Additional Resources



## User Stories

As a developer, I want to describe all of the dependencies in my project, so that they are all accounted for so that I don't need to manually hunt down all of the repos
* Describes a config file, dep.yaml

As a developer, I want to initiate/bootstrap my environment easily, so that I don't need to clone and setup each project/repo
* Describes the initial `clone` of all, along with potentially a `npm install` and/or `composer install` where neccessary

As a developer, I want to orchestrate actions across a subset of relevant repos so that a coordinated release can be prepared for multiple services.
* Describes a method to filter dependencies

As a developer, I want to define a pattern for subsets of projects to follow so that I can issue commands that act on all at once.
* This describes the the shortcoming of git submodules, you can't easily perform opperations on just a few.

As a developer, I want to allow fellow developers to also work within my project subset easily so that we don't need to recreate my setup everytime a new developer steps in.
* (Trying to describe some form of lock file)
* Lock files would cause so many merge conflicts if you wanted each branch to have it's own lock that tracked that branch
  * Built-in variables could fix a majority of this
  * Instead of saying something like `- branch: develop`, you could say something like `- branch: $BRANCH` which will look for the active branch, and use that.
     * This implies a parent concept, though. But the dep.yaml has to live somewhere.



## License

- [LICENSE](LICENSE) (Expat/[MIT License][MIT])

[MIT]: http://www.opensource.org/licenses/MIT "The MIT License (MIT)"