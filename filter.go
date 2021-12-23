package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/**
Flattens MergeDeps into a single-level array, and appends that onto depList
This is due to YAML limitations:
	// Step 0: Flatten merge-deps with deps. Because YAML doesn't support merging sequences:
	//  http://yaml4r.sourceforge.net/doc/page/collections_in_yaml.htm
	// 	https://stackoverflow.com/a/30770740/663058
*/
func flatten(deps []dep) []dep {
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
Flattens the dep YAML
	dep.Labels are inherited
	dep.Location is expanded
*/
func unwrap(deps []dep, baseDir string, labels []string, params map[string]string) []dep {
	var foundDeps []dep
	for _, dep := range deps {
		dep.Location = filepath.Clean(baseDir + dep.Location)
		dep.Labels = append(dep.Labels, labels...) // Inherit labels

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

		dep.Params = mergeMap(dep.Params, params) // Inherit parents Params
		if dep.DepList != nil {
			foundDeps = append(foundDeps, unwrap(dep.DepList, dep.Location+string(os.PathSeparator), dep.Labels, parentParms(dep.Params))...)
			dep.DepList = nil
		}

		foundDeps = append(foundDeps, dep)
	}

	return foundDeps
}

/*
In mustache.go the `func lookup()` attempts to expand dot notation, so I forked that codebase
and commented that out to allow directory style parental param lookups.
*/
func parentParms(params map[string]string) map[string]string {
	parents := make(map[string]string)
	for name, param := range params {
		parents[".."+string(os.PathSeparator)+name] = param
	}

	return parents
}

/**
Merge the parent map with the sibling, but allow sibling to take precedence if there are dupes.
*/
func mergeMap(sibling map[string]string, parent map[string]string) map[string]string {
	for name, param := range parent {
		if _, ok := sibling[name]; !ok {
			sibling[name] = param
		}
	}
	return sibling
}

/**
Filters out to just a kind
*/
func applyFilterKind(deps []dep, kind string) []dep {
	var foundDeps []dep
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
func applyFilterLabel(deps []dep, perform perform) []dep {
	if perform.Labels == "" {
		if perform.Verbose || !perform.Force {
			fmt.Println("Warning: No labels, using all deps of kind.")
		}

		if !perform.Force && !askForConfirmation("Are you sure you want to continue?\n"+
			"(use --force to suppress this prompt.)", nil) {
			if perform.Verbose {
				fmt.Println("DepCharge cancelled.")
			}
			os.Exit(0)
		}
		// If no labels, and only kind, return that.
		return deps
	}

	labels := strings.Split(perform.Labels, ",")

	var foundDeps []dep
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
