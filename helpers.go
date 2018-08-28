package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"strings"
)

var mockDefaultAction func(chan<- bool, Dep, Perform)
func depInjDefaultAction() func(chan<- bool, Dep, Perform) {
	if isTesting() {
		if mockDefaultAction != nil {
			return mockDefaultAction
		} else {
			return placeholderDefaultActionMock
		}
	}
	return defaultAction
}

func placeholderDefaultActionMock(complete chan<- bool, dep Dep, perform Perform) {
	fmt.Println("Fallback mockDefaultAction")
	fmt.Println("Should mock out for:")
	fmt.Println(dep, perform)
	// TODO: Setup a strict env var, and fail here
	//os.Exit(-1)
	complete <- true
}

func isTesting() bool {
	if flag.Lookup("test.v") == nil {
		return false
	}

	return true
}

func drainChannel(muchness int, toDrain <-chan bool){
	// In the case we run parallel, block until all goroutines signify completed.
	for i := 0; i < muchness; i++ {
		<-toDrain
	}
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

//func dumpStruct(depList []Dep) string {
//	fmt.Println("Dumping JSON:")
//	newYaml, _ := yaml.Marshal(depList)
//	newJson, _ := yaml.YAMLToJSON(newYaml)
//	fmt.Println(string(newJson))
//	return string(newJson)
//}

// https://gist.github.com/albrow/5882501
// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
// input gets an array from the user.
func askForConfirmation(request string, in *os.File) bool {
	if in == nil {
		in = os.Stdin
	}

	fmt.Println(request)
	fmt.Print("[y|N]: ")

	var response string
	_, err := fmt.Fscanln(in, &response)
	if err != nil {
		log.Fatal(err)
	}
	response = strings.ToLower(response)
	okayResponses := []string{"y", "yes"}
	nokayResponses := []string{"n", "no"}

	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation(request, in)
	}
}


// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}