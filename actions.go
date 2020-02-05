package main

import (
	"fmt"
	"github.com/centerorbit/mustache"
	"io"
	"os"
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

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		err := cmd.Start()
		if err != nil {
			if perform.Verbose {
				fmt.Println("Couldn't start command: ", perform.Kind+" "+strings.Join(mustachedActionParams, " "))
				fmt.Println("Due to this error:", err)
			}
		} else {
			go func() { _, _ = io.Copy(os.Stdout, stdout) }()
			go func() { _, _ = io.Copy(os.Stderr, stderr) }()
			err = cmd.Wait()

			if err != nil {
				if perform.Verbose {
					fmt.Println("Command finished with error: ", err)
					fmt.Println(perform.Kind + " " + strings.Join(mustachedActionParams, " "))
				}
			}
		}
	}

	complete <- true
}

/// ***  Helpers *** ///

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
