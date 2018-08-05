package main

import (
		"github.com/cbroglie/mustache"
		"os/exec"
	"fmt"
)

/// *** Actions *** ///

func defaultAction(complete chan<- bool, dep Dep, perform Perform) {

	mustachedActionParams := templateParams(dep, perform)

	fmt.Println("Running '", perform.Kind, mustachedActionParams, " for: ", dep.Location)
	var err error = nil
	var out []byte = nil

	if perform.DryRun {
		fmt.Println("Dry run of: ")
		fmt.Println(perform.Kind, " ", mustachedActionParams)
	} else {
		cmd := exec.Command(perform.Kind, mustachedActionParams...)
		cmd.Dir = dep.Location
		//TODO: Find a way to "stream" output to terminal?
		out, err = cmd.CombinedOutput() //Combines errors to output
		//out, err := cmd.Output() // just stdout
	}


	if err != nil {
		fmt.Println("Command finished with error: ", err)
	}

	if string(out) == "" {
		fmt.Println("Done!")
	} else {
		fmt.Println(string(out))
	}

	complete <- true
}

/// ***  Helpers *** ///

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
