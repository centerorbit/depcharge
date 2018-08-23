package main

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"os/exec"
	"strings"
)

/// *** Actions *** ///

func defaultAction(complete chan<- bool, dep Dep, perform Perform) {

	mustachedActionParams := templateParams(dep, perform)

	if perform.DryRun {
		fmt.Println("Dry run of: `", perform.Kind, strings.Join(mustachedActionParams, " "), "` for: ", dep.Location)
	} else {
		fmt.Println("Running: `", perform.Kind, strings.Join(mustachedActionParams, " "), "` for: ", dep.Location)
		command := perform.Kind + " " + strings.Join(mustachedActionParams, " ")
		cmd := exec.Command(perform.Kind, mustachedActionParams...)
		cmd.Dir = dep.Location
		//TODO: Find a way to "stream" output to terminal?
		out, err := cmd.CombinedOutput()
		checkOkay(command, out, err) //Combines errors to output
		//out, err := cmd.Output() // just stdout
	}

	complete <- true
}

func prepDockerComposeAction(dep Dep, perform Perform) string {

	mustachedActionParams := templateParams(dep, perform)

	return strings.Join(mustachedActionParams, " ")
}

func dockerComposeAction(complete chan<- bool, perform Perform, action []string) {
	if perform.DryRun {
		fmt.Println("Dry run of: ")
		fmt.Println(perform.Kind, action)
	} else {
		command := perform.Kind + " " + strings.Join(action, " ")
		cmd := exec.Command(perform.Kind, action...)
		//TODO: Find a way to "stream" output to terminal?
		out, err := cmd.CombinedOutput()
		checkOkay(command, out, err) //Combines errors to output
	}

	complete <- true
}

/// ***  Helpers *** ///

func checkOkay(command string, out []byte, err error) {
	if err != nil {
		fmt.Println("Command finished with error: ", err)
		fmt.Println(command)
	}

	if string(out) == "" {
		fmt.Println("Done!")
	} else {
		fmt.Println(string(out))
	}
}

func templateParams(dep Dep, perform Perform) []string {
	// Adding kind, name, and location to possible template params
	if dep.Params == nil {
		dep.Params = map[string]string{}
	}
	if _, ok := dep.Params["kind"]; !ok {
		dep.Params["kind"] = dep.Kind
	}
	if _, ok := dep.Params["name"]; !ok {
		dep.Params["name"] = dep.Name
	}
	if _, ok := dep.Params["location"]; !ok {
		dep.Params["location"] = dep.Location
	}

	mustachedActionParams := applyMustache(dep.Params, perform.Action)

	return mustachedActionParams
}

func applyMustache(params map[string]string, actionParams []string) []string {
	var mustachedActionParams []string
	mustache.AllowMissingVariables = false

	for _, value := range actionParams {
		data, err := mustache.Render(value, params)

		if err != nil {
			fmt.Println("Warning: ", err)
		}

		mustachedActionParams = append(mustachedActionParams, data)
	}

	return mustachedActionParams
}
