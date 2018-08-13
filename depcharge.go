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
	MergeDeps [][]Dep `json:"merge-deps"`
	DepList  []Dep    `json:"deps"`
	Labels	 []string `json:"labels"`
	Params map[string]string `json:"params"`
}

type DepList struct {
	Deps []Dep `json:"deps"`
}

type Perform struct {
	Kind string
	Instead string
	Labels string
	Action []string
	Exclusive bool
	DryRun bool
	Force bool
}


func main() {
	perform := processArgs()
	depList := load()


	// Figure out our current directory.
	cwd, _ := filepath.Abs("./")
	// ABS, and Clean strip trailing slash, we need to consistently re-add it for concating to work
	cwd += "/"

	// Step 0: Flatten merge-deps with deps. Because YAML doesn't support merging sequences:
	//  http://yaml4r.sourceforge.net/doc/page/collections_in_yaml.htm
	// 	https://stackoverflow.com/a/30770740/663058
	flattened := flatten(depList.Deps)
	//dumpStruct(flattened)

	// Step 1: We must flatten our YAML Struct, expanding location, and inheriting labels
	expanded := unwrap(flattened, cwd, nil)

	// Step 2: Now, lets filter out to only the kind that we want
	kindFiltered := applyFilterKind(expanded, perform.Kind)

	// Step 3: Filter out via labels, Label filtering is always inclusive.
	//  If a dep has the label (or its parent had it, hence inherited):
	//    It wins!
	labelFiltered := applyFilterLabel(kindFiltered, perform)

	// If '--instead' is provided, swap it out for Kind, _after_ filtering has been done
	if perform.Instead != "" {
		perform.Kind = perform.Instead;
	}

	// We select a handler based on our kind
	handler := findActionHandler(perform.Kind)


	// Finally, call the handler which will find and execute the kind+action
	//  across all final deps.
	complete := make(chan bool)

	n := handler(complete, labelFiltered, perform)

	for i := 0; i < n; i++ {
		<-complete
	}
	fmt.Println("depcharge complete!")
}

/// ***  Helpers *** ///

func dumpStruct(depList []Dep) string {
	fmt.Println("Dumping JSON:")
	newYaml, _ := yaml.Marshal(depList)
	newJson, _ := yaml.YAMLToJSON(newYaml)
	fmt.Println(string(newJson))
	return string(newJson)
}
