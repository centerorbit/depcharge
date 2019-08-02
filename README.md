# DepCharge
[![Go Report Card](https://goreportcard.com/badge/github.com/centerorbit/depcharge)](https://goreportcard.com/report/github.com/centerorbit/depcharge)
[![Build Status](https://cloud.drone.io/api/badges/centerorbit/depcharge/status.svg)](https://cloud.drone.io/centerorbit/depcharge)
[![coverage report](https://gitlab.com/centerorbit/depcharge/badges/master/coverage.svg)](https://gitlab.com/centerorbit/depcharge/commits/master)
[![GitHub license](https://img.shields.io/github/license/centerorbit/depcharge.svg)](https://github.com/centerorbit/depcharge/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/centerorbit/depcharge.svg)](https://github.com/centerorbit/depcharge/releases/latest)
[![Maintainability](https://api.codeclimate.com/v1/badges/5ef8ce4f942696ebace7/maintainability)](https://codeclimate.com/github/centerorbit/depcharge/maintainability)
<a href="https://github.com/avelino/awesome-go"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg"></a>


DepCharge is a tool that helps orchestrate the execution of commands across the many dependencies and directories in larger projects. It also provides an excellent mechanism for self-documenting a project's vast (and often disparate) dependencies.

**Sneak peek:**

	depcharge --labels=api -- git clone {{repo}} {{location}}
	
Will run `git clone` across all listed git dependencies with the label of "api" in your project (where submodules use to rule the land)


## Introduction

By creating a YAML file that describes all of your project(s) dependencies, you can then execute commands across all of them simultaneously.

A medium-to-large sized project (especially when using a microservice architecture) will consist of 3 or more separate repositories, and rely on a variety of package managers depending on the various languages chosen for each service. Typically, these repos must be managed, tracked, and released in some semblance of unison so that the dependant service calls can be understood and responded to appropriately.

For small (to single) teams, a single developer will often need to propagating and perform the same commands across all relevant services. This is a tedious, manual, and error-prone process that can occur every release.

DepCharge is designed to help fix that.

## How it Works & Usage
DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once. All of the examples here are just that: examples. DepCharge is designed to be as flexible as possible, so if you happen to use tools other than what's listed, they should work as well!

Usage: `depcharge [version] [--help|-h] [--dryrun|-d] [--force|-f] [--kind|-k=<kind>] [--instead|-x=<action>] [[--include|-i][-e|--exclude]] [--labels|-l=<comma-separated,inherited>] [--serial|-s] -- COMMAND [ARGS...]`

### Features:
* Built-in mustache templating, allows you to parametrize your commands
* Supports arbitrary params in YAML, whatever 'params: key: value' pairs you want (particularly useful for mustache templating)
* Supports YAML anchors
  * Even went the extra mile to support anchors + sequence merging via `merge-deps:` (see: `YAML Anchors & Sequences`)

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

### Flags

       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -d --dryrun  Will print out the command to be run, does not make changes to your system.
    -e --exclusive  Applies labels in an exclusive way (default).
    -f --force  Will force-run a command without confirmations, could be dangerous.
    -i --inclusive  Applies labels in an inclusive way.
    -k --kind  Targets specific kinds of dependencies (i.e. git, npm, composer)
    -l --labels  Filters to specific labels.
    -s --serial  Prevents parallel execution, runs commands one at a time.
    -x --instead  Instead of 'kind', perform a different command.


### Example commands:

Will run `git clone <repo> <location>` across all git dependencies:

	depcharge --kind=git -- clone {{repo}} {{location}}
	
Or, shorthand:

	depcharge -- git clone {{repo}} {{location}}
	
Will run `git status` across all git dependencies:

	depcharge -- git status
	
Will run `npm install` across any npm dependencies that have the label 'public':

	depcharge --labels=public -- npm install
	
Will run `composer install` across any composer dependencies that have either the label 'api', or 'soap':

	depcharge --inclusive --labels=api,soap -- composer install
	
And much more!

## YAML Anchors & Sequences
Due to a limitation in YAML itself, you cannot use anchors to merge sequences (arrays). Therefore this is programatically supported within DepCharge.

Invalid YAML, you cannot mix sequences `-` with anchors `*<name>` directly, this doesn't work:
```yaml
...
    deps:
      - kind: git
      - *composer
      - *vue
...
```
It's a beautiful concept though, that really helps with reusability and simplifies the overall YAML file, and so `merge-deps` was introduced to work around this shortcoming.

Working around this with merge-deps:
```yaml
.vue: &vue
  - name: Vue.js
    kind: npm
.composer: &composer
  - name: lumen
    kind: composer
    
deps:
  - name: ui
    kind: project
    location: ./code/app
    labels:
      - ui
    deps:
      - kind: git
        params:
          repo: git@example.com/ui.git
    merge-deps:
      - *composer
      - *vue
```

In the above example, `merge-deps:` supports listing your anchors, and these will then be expanded, then flattened and merged into `deps:` before final processing begins.


## Special Action Handlers
DepCharge has the ability to offer special-case action handlers. Specifically for situations where executing bulk commands cause difficulties and/or there are unexpected rough edges.

1. `git clone`
This is treated specially, in the sense that a regular clone will not act if parent directories aren't already in place. DepCharge will detect the `clone` action explicitly and attempt to create any parent directories before passing the command directly onto `git`

## Additional Resources
* https://mustache.github.io/

## Contributing
See: [CONTRIBUTING.md](CONTRIBUTING.md)

## License
- [LICENSE](LICENSE) (Expat/[MIT License][MIT])

[MIT]: http://www.opensource.org/licenses/MIT "The MIT License (MIT)"