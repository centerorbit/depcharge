package main

import (
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"os"
)

func depInjDockerComposeAction() func(complete chan<- bool, perform Perform, action []string) {
	if isTesting() {
		return func(complete chan<- bool, perform Perform, action []string) {
			fmt.Println("Mock dockerComposeAction")
		}
	}
	return dockerComposeAction
}

func depInjDefaultAction() func(chan<- bool, Dep, Perform) {
	if isTesting() {
		return func(complete chan<- bool, dep Dep, perform Perform) {
			fmt.Println("Mock dockerComposeAction")
		}
	}
	return defaultAction
}

func isTesting() bool {
	if flag.Lookup("test.v") == nil {
		return false
	}

	return true
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func dumpStruct(depList []Dep) string {
	fmt.Println("Dumping JSON:")
	newYaml, _ := yaml.Marshal(depList)
	newJson, _ := yaml.YAMLToJSON(newYaml)
	fmt.Println(string(newJson))
	return string(newJson)
}
