package main

import (
	"os"
	"fmt"
)

func findActionHandler(kind string) func([]Dep, Perform){
	switch kind {
	case "git":
		return gitActionHandler
	case "secret":
		return secretesActionHandler
	default:
		return defaultActionHandler
	}
}


/// *** Action Handlers *** ///

func defaultActionHandler(deps []Dep, perform Perform){
	for _, dep := range deps {
		defaultAction(dep, perform)
	}
}

func gitActionHandler(deps []Dep, perform Perform){
	for _, dep := range deps {
		switch perform.Action[0] {
		case "clone": // Clone breaks if the parent dirs aren't already there.
			there, _ := exists(dep.Location)
			if !there {
				if perform.DryRun {
					fmt.Println("DryRun, would have performed a: ", "`mkdir -p ", dep.Location, "`")
				} else {
					// mkdir -p <location>
					os.MkdirAll(dep.Location, os.ModePerm)
				}
			}

			defaultAction(dep, perform)

		default:
			defaultAction(dep, perform)
		}
	}
}


// TODO: make a special handler for secrets
func secretesActionHandler(deps []Dep, perform Perform){
	for _, dep := range deps {
		switch perform.Action[0] {
		//case "get": // Or something along these lines
		// Right now, everything falls through to default.
		default:
			defaultAction(dep, perform)
		}
	}
}


/// ***  Helpers *** ///

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}