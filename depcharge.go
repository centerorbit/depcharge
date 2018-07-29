package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/_examples"
	"os"
	"flag"
	"path/filepath"
	"os/exec"
	"strings"
	"time"
)

type Dep struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Location string   `json:"location"`
	Repo     string   `json:"repo"`
	DepList  []Dep    `json:"deps"`
	Labels	 []string `json:"labels"`
}

type DepList struct {
	Deps []Dep `json:"deps"`
}


/**
dep --kind=git --label=service,api status
dep --kind=npm install
 */

func main() {
	// Define, grab, and parse our args.
	kindPtr := flag.String("kind", "", "Targets specific kinds of dependencies (i.e. git, npm, composer)")
	labelPtr := flag.String("labels", "", "Filters to specific labels.")
	flag.Parse()
	action := flag.Args()

	// Read in our YAML file.
	yamlFile, err := ioutil.ReadFile("dep.yml")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// Unmarshal the YAML into a struct.
	var depList DepList
	err = yaml.Unmarshal(yamlFile, &depList)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// Figure out our current directory.
	cwd, _ := filepath.Abs("./")
	// ABS, and Clean strip trailing slash, we need to consistently re-add it for concating to work
	cwd += "/"

	// Step 1: We must flatten our YAML Struct, expanding location, and inheriting labels
	expanded := unwrap(depList.Deps, cwd, nil)

	// Step 2: Now, lets filter out to only the kind that we want
	kindFiltered := applyFilterKind(expanded, *kindPtr)

	// Step 3: Filter out via labels, Label filtering is always inclusive.
	//  If a dep has the label (or its parent had it, hence inherited):
	//    It wins!
	labelFiltered := applyFilterLabel(kindFiltered, *labelPtr)

	// Debugging, will output JSON of final filtered down deps
	//dumpStruct(labelFiltered)

	// We select a handler based on our kind
	handler := findActionHandler(*kindPtr)

	// Finally, call the handler which will find and execute the kind+action
	//  across all final deps.
	handler(*kindPtr, action, labelFiltered)
}

/// *** Primary Methods *** ///


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
func applyFilterLabel(deps []Dep, labelString string) []Dep {
	// If no labels, and onlyo kind, return that.
	if labelString == "" {
		fmt.Println("Warning: No labels, using all deps of kind.")
		fmt.Println("Press Ctrl+c to cancel...")
		time.Sleep(5 * time.Second)
		return deps
	}

	labels := strings.Split(labelString, ",")

	var foundDeps []Dep
	for _, dep := range deps {
		//Filter to || of labels
		for _, depLabel := range dep.Labels {
			for _, filterLabel := range labels {
				if filterLabel == depLabel {
					fmt.Println("Found a match!", dep)
					foundDeps = append(foundDeps, dep)
				}
			}
		}
	}

	return foundDeps
}


func findActionHandler(kind string) func(string, []string, []Dep){
	switch kind {
	case "git":
		return gitActionHandler
	default:
		return defaultActionHandler
	}
}


/// *** Action Handlers *** ///

func defaultActionHandler(kind string, action []string, deps []Dep){
	for _, dep := range deps {
		defaultAction(kind, action, dep)
	}
}

func gitActionHandler(kind string, action []string, deps []Dep){
	for _, dep := range deps {
		switch action[0] {
		// TODO: remove this once params are working
		// need to use "repo" from yaml for clone, this differs per dep.
		case "clone":
			fmt.Println("Attempting to clone " + dep.Repo + " to: " + dep.Location)
			there, _ := exists(dep.Location)
			if there {
				fmt.Println("Dir exists, fetching instead.")
				defaultAction(kind, []string{"fetch"}, dep)
			} else {
				gitClone(dep.Location, dep.Repo)
			}
		// Right now, everything falls through to default.
		default:
			defaultAction(kind, action, dep)
		}
	}
}


/// *** Actions *** ///

func defaultAction(kind string, action []string, dep Dep) {
	fmt.Println("Running '", kind, strings.Join(action, " "), "' for: ", dep.Location)

	cmd := exec.Command(kind, action...)
	cmd.Dir = dep.Location
	//TODO: Find a way to "stream" output to terminal?
	out, err := cmd.CombinedOutput() //Combines errors to output
	//out, err := cmd.Output() // just stdout

	if err != nil {
		fmt.Println("Command finished with error: ", err)
	}

	if string(out) == "" {
		fmt.Println("Done!")
	} else {
		fmt.Println(string(out))
	}
}


func gitClone(path string, url string) {
	// Clone the given repository to the given directory
	examples.Info("git clone "+url+" "+path)

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	examples.CheckIfError(err)
}


/// ***  Helpers *** ///

func dumpStruct(depList []Dep) string {
	fmt.Println("Dumping JSON:")
	newYaml, _ := yaml.Marshal(depList)
	newJson, _ := yaml.YAMLToJSON(newYaml)
	fmt.Println(string(newJson))
	return string(newJson)
}


// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}