package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/integrii/flaggy"
	"io/ioutil"
	"os"
)

var version string

func processArgs() Perform {
	flaggy.SetVersion(version)

	var perform Perform

	// Define, grab, and parse our args.
	kind := ""
	flaggy.String(&kind, "k", "kind", "Targets specific kinds of dependencies (i.e. git, npm, composer)")

	instead := ""
	flaggy.String(&instead, "x", "instead", "Instead of 'kind', perform a different command.")

	labels := ""
	flaggy.String(&labels, "l", "labels", "Filters to specific labels.")

	exclusive := false
	flaggy.Bool(&exclusive, "e", "exclusive", "Applies labels in an exclusive way (default).")

	inclusive := false
	flaggy.Bool(&inclusive, "i", "inclusive", "Applies labels in an inclusive way.")

	dryRun := false
	flaggy.Bool(&dryRun, "d", "dryrun", "Will print out the command to be run, does not make changes to your system.")

	flaggy.SetDescription(" a tool designed to help orchestrate the execution of commands across many directories at once.")

	flaggy.DefaultParser.AdditionalHelpPrepend = "\n" +
		"Use -- to separate depcharge commands from intended execution commands."

	flaggy.DefaultParser.AdditionalHelpAppend =
		"Usage: depcharge [--kind=<kind>] [--instead=<action>] [--labels=<comma-separated,inherited>] [OPTIONS...] -- [COMMAND/ARGS...]" +
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
			"" +
			"\n\n" +
			"Primary Commands: \n" +
			"--kind		Is the top-level filter that's applied, opperations are run based on 'kind' \n" +
			"			if --kind is not specified, then the first COMMAND/ARG is used \n" +
			"--instead	Is used to specify a command you'd like to run against --kind, but is not 'kind'. \n" +
			"--labels	Comma separated list of labels to filter by, inherited from parents \n" +
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

	flaggy.Parse()
	action := flaggy.TrailingArguments

	if exclusive && inclusive {
		fmt.Println("--exclusive and --inclusive cannot be specified at the same time.")
		os.Exit(-1)
	}

	exclusive = !inclusive

	if kind == "" && len(action) == 0 {
		flaggy.ShowHelpAndExit("\n ERROR: You must provide at least a '--kind' or one ARG.")
	}
	if kind == "" && len(action) >= 1 {
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
	perform.DryRun = dryRun

	return perform
}

func load() DepList {
	// Read in our YAML file.
	yamlFile, err := ioutil.ReadFile("dep.yml")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(-1)
	}

	// Unmarshal the YAML into a struct.
	var depList DepList
	err = yaml.Unmarshal(yamlFile, &depList)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(-1)
	}

	return depList
}
