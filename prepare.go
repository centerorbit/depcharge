package main

import (
	"os"
	"fmt"
	"flag"
	"io/ioutil"
	"github.com/ghodss/yaml"
)


func processArgs() Perform {
	var perform Perform

	// Define, grab, and parse our args.
	kindPtr := flag.String("kind", "", "Targets specific kinds of dependencies (i.e. git, npm, composer)")
	labelPtr := flag.String("labels", "", "Filters to specific labels.")
	exclusiveFlag := flag.Bool("exclusive", false, "Applies labels in an exclusive way (default).")
	inclusiveFlag := flag.Bool("inclusive", false, "Applies labels in an inclusive way.")
	dryRunFlag := flag.Bool("dryrun", false, "Will print out the command to be run, does not make changes to your system.")
	helpFlag := flag.Bool("help", false, "Prints the help text.")

	flag.Parse()
	action := flag.Args()

	if *helpFlag {
		fmt.Println(
			"\n" +
				"DepCharge is a tool designed to help orchestrate the execution of commands across many directories at once." +
				"\n\n" +
				"Usage: depcharge --kind=<kind> [--labels=<comma-separated,inherited>] [OPTIONS...] COMMAND [ARGS...]" +
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
				" --kind		Is the top-level filter that's applied, opperations are run based on 'kind' \n" +
				" --labels		Comma separated list of labels to filter by, inherited from parents \n" +
				"\n" +
				"Available Options: \n" +
				" --help		Shows this message \n" +
				" --dryrun		Prints out intended command without executing it \n" +
				" --exclusive	(default) For a match to be found, it must contain at least all provided labels \n" +
				" --inclusive   For a match to be found, it must contain at least one of the provided labels \n" +
				"\n" +
				"Example commands:" +
				"\n" +
				"Will run `git clone <location>` across all git dependencies: \n" +
				"	depcharge --kind=git clone {{location}}" +
				"\n\n" +
				"Will run `git status` across all git dependencies: \n" +
				"	depcharge --kind=git status" +
				"\n\n" +
				"Will run `npm install` across any npm dependencies that have the label 'public': \n" +
				"	depcharge --kind=npm --labels=public install" +
				"\n\n" +
				"Will run `composer install` across any composer dependencies that have either the label 'api', or 'soap': \n" +
				"	depcharge --kind=composer --inclusive --labels=api,soap install" +
				"")
		os.Exit(0)
	}

	if *exclusiveFlag && *inclusiveFlag {
		fmt.Println("--exclusive and --inclusive cannot be specified at the same time.")
		os.Exit(-1)
	}

	exclusive := !*inclusiveFlag

	perform.Kind = *kindPtr
	perform.Labels = *labelPtr
	perform.Action = action
	perform.Exclusive = exclusive
	perform.DryRun = *dryRunFlag

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