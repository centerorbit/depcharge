package main

import (
	"fmt"
	"github.com/centerorbit/mustache"
	"os/exec"
	"strings"
)

/// *** Actions *** ///

var execCommand = exec.Command

func defaultAction(complete chan<- bool, dep dep, perform perform) {

	mustachedActionParams := applyMustache(dep.Params, perform.Action, perform.Verbose)

	if perform.DryRun && perform.Verbose {
		fmt.Println("Dry run of: `", perform.Kind, strings.Join(mustachedActionParams, " "), "` for: ", dep.Location)
	} else {
		if perform.Verbose {
			fmt.Println("Running: `", perform.Kind, strings.Join(mustachedActionParams, " "), "` for: ", dep.Location)
		}
		cmd := execCommand(perform.Kind, mustachedActionParams...)
		cmd.Dir = dep.Location
		//TODO: Find a way to "stream" output to terminal?
		out, err := cmd.CombinedOutput()

		command := perform.Kind + " " + strings.Join(mustachedActionParams, " ")
		checkOkay(command, out, err, perform.Verbose) //Combines errors to output
		//out, err := cmd.Output() // just stdout
	}

	complete <- true
}

/// ***  Helpers *** ///
var checkOkayIntercept func(command string, out []byte, err error)

func checkOkay(command string, out []byte, err error, verbose bool) {
	if checkOkayIntercept != nil {
		checkOkayIntercept(command, out, err)
		return
	}

	if err != nil {
		if verbose {
			fmt.Println("Command finished with error: ", err)
			fmt.Println(command)
		} else {
			fmt.Println(err)
		}
	}

	if string(out) == "" {
		if verbose {
			fmt.Println("Done!")
		}
	} else {
		fmt.Print(string(out))
	}
}

func applyMustache(params map[string]string, actionParams []string, verbose bool) []string {
	var mustachedActionParams []string
	mustache.AllowMissingVariables = false

	for _, value := range actionParams {

		data, err := mustache.Render(value, params)

		if err != nil && verbose {
			fmt.Println("Warning: ", err)
		}

		mustachedActionParams = append(mustachedActionParams, data)
	}

	return mustachedActionParams
}
