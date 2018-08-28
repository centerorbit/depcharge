package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
)


func fakeExecCommand(command string, args...string) *exec.Cmd {
	fmt.Println("Hit fake Exec with", command, args)
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestDefaultAction(t *testing.T) {
	fmt.Println("Begin testing default action")
	execCommand = fakeExecCommand
	defer func(){
		execCommand = exec.Command
		checkOkayIntercept = nil
	}()

	checkOkayIntercept = func(command string, out []byte, err error){
		assert.Equal(t, "git status", command)
		assert.Nil(t, err)
		fmt.Println(string(out))
		assert.Equal(t, "TestHelperProcess, intercepting exec.Command", string(out))
	}

	complete := make(chan bool, 1)

	dep := Dep{
		Kind:     "git",
		Name:     "depcharge",
		Location: "./",
	}

	perform := Perform{
		Kind: "git",
		Action: []string{"status"},
	}

	defaultAction(complete, dep, perform)
	<-complete

	fmt.Println("Done testing default action")

}

func TestHelperProcess(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprint(os.Stdout, "TestHelperProcess, intercepting exec.Command")
	// some code here to check arguments perhaps?
	//fmt.Fprint(os.Stdout, os.Args)
	os.Exit(0)
}


func TestTemplateParams(t *testing.T) {

	dep := Dep{
		Kind:     "git",
		Name:     "depcharge",
		Location: "./",
	}

	perform := Perform{
		Action: []string{
			"To {{kind}}, or",
			"{{name}},",
			"that is the {{answer}}.",
		},
	}

	results := templateParams(dep, perform)
	assert.Equal(t,
		"To git, or depcharge, that is the ",
		strings.Join(results, " "))


	dep.Params = map[string]string{
		"is":     "be",
		"is not": "not to be",
		"answer": "question",
	}
	results = templateParams(dep, perform)
	assert.Equal(t,
		"To git, or depcharge, that is the question.",
		strings.Join(results, " "))
}

func TestApplyMustache(t *testing.T) {
	params := map[string]string{
		"is":     "be",
		"is not": "not to be",
		"answer": "question",
	}

	actionParams := []string{
		"To {{is}}, or",
		"{{is not}},",
		"that is the {{answer}}.",
	}

	results := applyMustache(params, actionParams)

	assert.Equal(t,
		"To be, or not to be, that is the question.",
		strings.Join(results, " "))

}

func TestCheckOkay(t *testing.T) {
	err := errors.New("Fake error")
	checkOkay("Not okay", nil, err)

	checkOkay("Is okay", nil, nil)

	out := []byte("Here is a string....")
	checkOkay("Is okay", out, nil)
}
