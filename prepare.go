package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/integrii/flaggy"
	"io/ioutil"
	"os"
)

var version string

const helpText = "Usage: depcharge [--labels=<comma-separated,inherited>] [OPTIONS...] -- [COMMAND/ARGS...]" +
	"\n\n" +
	"Features:" +
	"\n" +
	"* Supports arbitrary params, whatever 'params: key: value' pairs you want \n" +
	"* Built-in mustache templating, allows you to parametrize your commands \n" +
	"* Supports YAML anchors \n" +
	"\n" +
	"Description:" +
	"\n" +
	"depcharge will read the dep.yml file in the current working directory, and \n" +
	"perform all commands relative to that location." +
	"\n\n" +
	"Example dep.yml:" +
	"\n" +
	"deps: \n" +
	"    - name: frontend \n" +
	"      kind: git \n" +
	"      location: ./app/frontend \n" +
	"      labels: \n" +
	"        - public \n" +
	"      params: \n" +
	"        repo: git@example.com:frontend.git \n" +
	"      deps: \n" +
	"        - name: vue.js \n" +
	"          kind: npm \n" +
	"    - name: backend \n" +
	"      kind: git \n" +
	"      location: ./app/backend \n" +
	"      labels: \n" +
	"        - api \n" +
	"      params: \n" +
	"        repo: git@example.com:backend.git \n" +
	"      deps: \n" +
	"        - name: lumen \n" +
	"          kind: composer \n" +
	"\n" +
	"Example commands:" +
	"\n" +
	"Will run `git clone <location>` across all git dependencies: \n" +
	"	depcharge --kind=git -- clone {{location}}" +
	" (same as:)	depcharge -- git clone {{location}}" +
	"\n\n" +
	"Will run `git status` across all git dependencies: \n" +
	"	depcharge -- git status" +
	"\n\n" +
	"Will run `npm install` across any npm dependencies that have the label 'public': \n" +
	"	depcharge --labels=public -- npm install" +
	"\n\n" +
	"Will run `composer install` across any composer dependencies that have either the label 'api', or 'soap': \n" +
	"	depcharge --inclusive --labels=api,soap -- composer install"

func processArgs() perform {
	flaggy.SetVersion(version)

	var perform perform

	// Define, grab, and parse our args.
	dryRun := false
	flaggy.Bool(&dryRun, "d", "dryrun", "Will print out the command to be run, does not make changes to your system.")

	exclusive := false
	flaggy.Bool(&exclusive, "e", "exclusive", "Applies labels in an exclusive way (default).")

	force := false
	flaggy.Bool(&force, "f", "force", "Will force-run a command without confirmations, could be dangerous.")

	inclusive := false
	flaggy.Bool(&inclusive, "i", "inclusive", "Applies labels in an inclusive way.")

	kind := ""
	flaggy.String(&kind, "k", "kind", "Targets specific kinds of dependencies (i.e. git, npm, composer)")

	labels := ""
	flaggy.String(&labels, "l", "labels", "Filters to specific labels.")

	serial := false
	flaggy.Bool(&serial, "s", "serial", "Prevents parallel execution, runs commands one at a time.")

	instead := ""
	flaggy.String(&instead, "x", "instead", "Instead of 'kind', perform a different command.")

	verbose := false
	flaggy.Bool(&verbose, "v", "verbose", "Will print out additional information.")

	flaggy.SetDescription(" a tool designed to help orchestrate the execution of commands across many directories at once.")

	flaggy.DefaultParser.AdditionalHelpPrepend = "\n" +
		"Use -- to separate depcharge commands from intended execution commands."

	flaggy.DefaultParser.AdditionalHelpAppend = helpText

	flaggy.Parse()
	action := flaggy.TrailingArguments

	if exclusive && inclusive {
		flaggy.ShowHelpAndExit("--exclusive and --inclusive cannot be specified at the same time.")
	}

	exclusive = !inclusive

	if kind == "" && len(action) == 0 {
		flaggy.ShowHelpAndExit("\n ERROR: You must provide at least a '--kind' or one ARG.")
	}

	if len(action) >= 1 && kind == action[0] {
		// If kind is the same as action, just use it once
		perform.Kind = kind
		perform.Action = action[1:]
	} else if kind == "" && len(action) >= 1 {
		// If kind isn't set, look at the first arg, and use that instead

		// First, grab the action (ie up, or build)
		perform.Kind = action[0]
		//Trim off the first action, so that it doesn't get loop-added by the deps
		perform.Action = action[1:]
	} else {
		perform.Kind = kind
		perform.Action = action
	}

	perform.Instead = instead
	perform.Labels = labels
	perform.Exclusive = exclusive
	perform.Serial = serial
	perform.DryRun = dryRun
	perform.Force = force
	perform.Verbose = verbose

	return perform
}

func load() depList {
	// Read in our YAML file.
	yamlFile, err := ioutil.ReadFile("dep.yml")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(-1)
	}

	// Unmarshal the YAML into a struct.
	var depList depList
	err = yaml.Unmarshal(yamlFile, &depList)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(-1)
	}

	return depList
}
