package main

import (
	"fmt"
	"github.com/ghodss/yaml"
						"path/filepath"
					)

type Dep struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Location string   `json:"location"`
	DepList  []Dep    `json:"deps"`
	Labels	 []string `json:"labels"`
	Params map[string]string `json:"params"`
}

type DepList struct {
	Deps []Dep `json:"deps"`
}

type Perform struct {
	Kind string
	Labels string
	Action []string
	Exclusive bool
	DryRun bool
}


func main() {
	perform := processArgs()
	depList := load()


	// Figure out our current directory.
	cwd, _ := filepath.Abs("./")
	// ABS, and Clean strip trailing slash, we need to consistently re-add it for concating to work
	cwd += "/"

	// Step 1: We must flatten our YAML Struct, expanding location, and inheriting labels
	expanded := unwrap(depList.Deps, cwd, nil)

	// Step 2: Now, lets filter out to only the kind that we want
	kindFiltered := applyFilterKind(expanded, perform.Kind)

	// Step 3: Filter out via labels, Label filtering is always inclusive.
	//  If a dep has the label (or its parent had it, hence inherited):
	//    It wins!
	labelFiltered := applyFilterLabel(kindFiltered, perform)

	// Debugging, will output JSON of final filtered down deps
	//dumpStruct(labelFiltered)

	// We select a handler based on our kind
	handler := findActionHandler(perform.Kind)

	// Finally, call the handler which will find and execute the kind+action
	//  across all final deps.
	handler(labelFiltered, perform)
}

/// ***  Helpers *** ///

func dumpStruct(depList []Dep) string {
	fmt.Println("Dumping JSON:")
	newYaml, _ := yaml.Marshal(depList)
	newJson, _ := yaml.YAMLToJSON(newYaml)
	fmt.Println(string(newJson))
	return string(newJson)
}
