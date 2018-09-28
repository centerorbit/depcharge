package main

import (
	"fmt"
	"os"
)

func findActionHandler(kind string) func(chan bool, []dep, perform) int {
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

func defaultActionHandler(complete chan bool, deps []dep, perform perform) int {
	performAction := depInjDefaultAction()

	n := 0
	for _, dep := range deps {
		n++
		go performAction(complete, dep, perform)

		// If Serial, block until complete on each.
		if perform.Serial {
			<-complete
		}
	}
	return n
}

func gitActionHandler(complete chan bool, deps []dep, perform perform) int {
	performAction := depInjDefaultAction()

	n := 0
	for _, dep := range deps {
		n++
		switch perform.Action[0] {
		case "clone": // Clone breaks if the parent dirs aren't already there.
			there, _ := exists(dep.Location)
			if !there {
				if perform.DryRun {
					fmt.Println("DryRun, would have performed a: `mkdir -p ", dep.Location, "`")
				} else {
					// mkdir -p <location>
					os.MkdirAll(dep.Location, os.ModePerm)
				}
			}

			go performAction(complete, dep, perform)

		default:
			go performAction(complete, dep, perform)
		}

		// If Serial, block until complete on each.
		if perform.Serial {
			<-complete
		}
	}

	return n
}

// TODO: make a special handler for secrets
func secretesActionHandler(complete chan bool, deps []dep, perform perform) int {
	performAction := depInjDefaultAction()

	n := 0
	for _, dep := range deps {
		n++
		switch perform.Action[0] {
		//case "get": // Or something along these lines
		// Right now, everything falls through to default.
		default:
			go performAction(complete, dep, perform)
		}

		// If Serial, block until complete on each.
		if perform.Serial {
			<-complete
		}
	}
	return n
}
