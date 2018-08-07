package main

import (
		"github.com/cbroglie/mustache"
		"os/exec"
	"fmt"
	"strings"
)

/// *** Actions *** ///

func defaultAction(complete chan<- bool, dep Dep, perform Perform) {

	mustachedActionParams := templateParams(dep, perform)

	fmt.Println("Running '", perform.Kind, mustachedActionParams, " for: ", dep.Location)

	if perform.DryRun {
		fmt.Println("Dry run of: ")
		fmt.Println(perform.Kind, " ", mustachedActionParams)
	} else {
		cmd := exec.Command(perform.Kind, mustachedActionParams...)
		cmd.Dir = dep.Location
		//TODO: Find a way to "stream" output to terminal?
		checkOkay(cmd.CombinedOutput()) //Combines errors to output
		//out, err := cmd.Output() // just stdout
	}

	complete <- true
}

func prepDockerComposeAction( dep Dep, perform Perform) string {

	mustachedActionParams := templateParams(dep, perform)

	return strings.Join(mustachedActionParams, " ")
}

func dockerComposeAction(complete chan<- bool, perform Perform, action []string)  {
	if perform.DryRun {
		fmt.Println("Dry run of: ")
		fmt.Println(perform.Kind, action)
	} else {
		cmd := exec.Command(perform.Kind, action...)
		//TODO: Find a way to "stream" output to terminal?
		//TODO: move checkOkay to better helpers location
		checkOkay(cmd.CombinedOutput()) //Combines errors to output
	}

	complete <- true
}

/// ***  Helpers *** ///

func checkOkay(out []byte, err error)  {
	if err != nil {
		fmt.Println("Command finished with error: ", err)
	}

	if string(out) == "" {
		fmt.Println("Done!")
	} else {
		fmt.Println(string(out))
	}
}

func templateParams(dep Dep, perform Perform) []string  {
	// Adding kind, name, and location to possible template params
	if dep.Params == nil { dep.Params = map[string]string{} }
	if _, ok := dep.Params["kind"];     !ok { dep.Params["kind"]     = dep.Kind	    }
	if _, ok := dep.Params["name"];     !ok { dep.Params["name"]     = dep.Name	    }
	if _, ok := dep.Params["location"]; !ok { dep.Params["location"] = dep.Location	}

	mustachedActionParams := applyMustache(dep.Params, perform.Action)

	return mustachedActionParams
}

func applyMustache (params map[string]string, actionParams []string) []string  {
	var mustachedActionParams []string

	for _, value := range(actionParams){
		data, _ := mustache.Render(value, params)
		mustachedActionParams = append(mustachedActionParams, data)
	}

	return mustachedActionParams
}
