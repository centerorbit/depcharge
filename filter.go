package main

import (
	"path/filepath"
	"fmt"
	"time"
	"strings"
)

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
			foundDeps = append(foundDeps, unwrap(dep.DepList, dep.Location + "/", dep.Labels)...)
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
		if ! perform.Force {
			fmt.Println("Press Ctrl+c to cancel...")
			time.Sleep(5 * time.Second)
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
				counter ++
			}
		}
	}
	if len(against) == counter{
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

