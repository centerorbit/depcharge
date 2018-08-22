package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"os"
)

/**
Flattens MergeDeps into a single-level array, and appends that onto DepList
This is due to YAML limitations:
	// Step 0: Flatten merge-deps with deps. Because YAML doesn't support merging sequences:
	//  http://yaml4r.sourceforge.net/doc/page/collections_in_yaml.htm
	// 	https://stackoverflow.com/a/30770740/663058
*/
func flatten(deps []Dep) []Dep {
	//Go through all of the deps, and check if they need flattening.
	for key, dep := range deps {

		//Cycle through the arrays of arrays of deps
		for _, mdep := range dep.MergeDeps {
			// and spread/merge those into the deps.
			dep.DepList = append(dep.DepList, flatten(mdep)...)
		}
		dep.MergeDeps = nil

		// The recursive part of this function
		// Called after the MergeDeps, in case the anchored deps have deps
		dep.DepList = flatten(dep.DepList)

		deps[key] = dep
	}

	return deps
}

/**
Flattens the Dep YAML
	dep.Labels are inherited
	dep.Location is expanded
*/
func unwrap(deps []Dep, baseDir string, labels []string) []Dep {
	var foundDeps []Dep
	for _, dep := range deps {
		dep.Location = filepath.Clean(baseDir + dep.Location)
		dep.Labels = append(dep.Labels, labels...) // Inherit labels
		if dep.DepList != nil {
			foundDeps = append(foundDeps, unwrap(dep.DepList, dep.Location+"/", dep.Labels)...)
			dep.DepList = nil
		}

		foundDeps = append(foundDeps, dep)
	}

	return foundDeps
}

/**
Filters out to just a kind
*/
func applyFilterKind(deps []Dep, kind string) []Dep {
	var foundDeps []Dep
	for _, dep := range deps {
		if dep.Kind != "" && dep.Kind == kind {
			foundDeps = append(foundDeps, dep)
		}
	}

	return foundDeps
}

/**
Applies filters
Splits comma separated
*/
func applyFilterLabel(deps []Dep, perform Perform) []Dep {
	if perform.Labels == "" {
		fmt.Println("Warning: No labels, using all deps of kind.")

		if !perform.Force && !askForConfirmation("Are you sure you want to continue?"){
			fmt.Println("DepCharge cancelled.")
			os.Exit(0)
		}
		// If no labels, and only kind, return that.
		return deps
	}

	labels := strings.Split(perform.Labels, ",")

	var foundDeps []Dep
	for _, dep := range deps { // Cycle through all of the deps

		var match bool
		if perform.Exclusive {
			match = isExclusive(dep.Labels, labels)
		} else {
			match = isInclusive(dep.Labels, labels)
		}

		if match {
			foundDeps = append(foundDeps, dep)
		}
	}

	return foundDeps
}

func isExclusive(what []string, against []string) bool {
	counter := 0
	for _, item := range what {
		for _, compare := range against {
			if item == compare {
				counter++
			}
		}
	}
	if len(against) == counter {
		return true
	}
	return false
}

func isInclusive(what []string, against []string) bool {
	if (len(what) == 0 && len(against) == 0) || len(against) == 0 {
		return true
	}

	match := false
InclusiveSearch:
	for _, item := range what {
		for _, compare := range against {
			if item == compare {
				match = true
				break InclusiveSearch
			}
		}
	}

	return match
}
