# DepCharge
[![pipeline status](https://gitlab.com/centerorbit/depcharge/badges/master/pipeline.svg)](https://gitlab.com/centerorbit/depcharge/commits/master) [![coverage report](https://gitlab.com/centerorbit/depcharge/badges/master/coverage.svg)](https://gitlab.com/centerorbit/depcharge/commits/master)

DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once. It also proves to provide an excellent mechanism for self-documenting a project's vast (and often disparate) dependencies.

**Sneak peek:**

	depcharge --labels=api -- git clone {{repo}} {{location}}
	
Will run `git clone` across all listed git dependencies with the label of "api" in your project (where submodules use to rule the land)


## Introduction
A medium-to-large sized project (especially when using a microservice architecture) will consist of 3 or more separate repositories, and rely on a variety of package managers depending on the various languages chosen for each service. Typically, these repos must be managed, tracked, and released in some semblance of unison so that the dependant service calls can be understood and responded to appropriately.

For small (to single) teams, a single developer will often need to propagating and perform the same commands across all relevant services. This is a tedious, manual, and error-prone process that can occur every release.

DepCharge is designed to help fix that.

By creating a YAML file that describes all of your project(s) dependencies, you can then execute commands across all of them simultaneously.

## How it Works & Usage
All of the examples here are just that: examples. DepCharge is designed to be as flexible as possible, so if you happen to use tools other than what's listed, they should work as well!

DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once.

Usage: `depcharge [--kind=<kind>] [--instead=<action>] [--labels=<comma-separated,inherited>] [OPTIONS...] -- COMMAND [ARGS...]`

### Features:
* Supports arbitrary params, whatever 'params: key: value' pairs you want
* Built-in mustache templating, allows you to parametrize your commands
* Supports YAML anchors
  * Even went the extra mile to support anchors + sequence merging via `MergeDeps`

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

 --kind		Is the top-level filter that's applied, opperations are run based on 'kind'. If --kind is not specified, then the first COMMAND/ARG is used
 --instead  Is used to specify a command you'd like to run against --kind, but is not 'kind'.
 --labels	Comma separated list of labels to filter by, inherited from parents

### Available Options:

 --help			Shows this message

 --dryrun		Prints out intended command without executing it

 --exclusive	(default) For a match to be found, it must contain at least all provided labels

 --inclusive   	For a match to be found, it must contain at least one of the provided labels

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
1. `docker-compose` This command is still in experimental development, but allows for a depcharge to specify a set of override files that would be passed to the `docker-compose` comand via `-f` flags to start whole subsets or entire local clusters of containers.
   * It currently does not support flags, only supports the base commands (up, build, down, etc.), and the args are flipped. Example:
     * Regular: `docker-compose -f file.yml up`
     * DepCharge: `depcharge --kind=docker-compose up {{location}}/{{file}}`
   * This is all subject to change as things get tested, and we near **v1.0**

## Additional Resources
* https://mustache.github.io/

## Contributing
* [GitLab is where development & pipelines occur](https://gitlab.com/centerorbit/depcharge)
* [GitHub is the public-access repository](https://github.com/centerorbit/depcharge)


## License
- [LICENSE](LICENSE) (Expat/[MIT License][MIT])

[MIT]: http://www.opensource.org/licenses/MIT "The MIT License (MIT)"