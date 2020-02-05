package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Thanks to: https://npf.io/2015/06/testing-exec-command/
//    and: https://github.com/golang/go/blob/master/src/os/exec/exec_test.go#L31
func fakeExecCommand(command string, args ...string) *exec.Cmd {
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
	defer func() {
		execCommand = exec.Command
	}()

	complete := make(chan bool, 1)

	dep := dep{
		Kind:     "git",
		Name:     "depcharge",
		Location: "./",
	}

	perform := perform{
		Kind:   "git",
		Action: []string{"status"},
	}

	defaultAction(complete, dep, perform)
	<-complete

	fmt.Println("Done testing default action")

}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	fmt.Fprint(os.Stdout, "TestHelperProcess, intercepting exec.Command")
	// some code here to check arguments perhaps?
	//fmt.Fprint(os.Stdout, os.Args)
	os.Exit(0)
}

func TestTemplateParams(t *testing.T) {

	dep := dep{
		Params: map[string]string{
			"kind":     "git",
			"name":     "depcharge",
			"location": "./",
			"answer":   "question",
		},
	}

	perform := perform{
		Action: []string{
			"To {{kind}}, or",
			"{{name}},",
			"that is the {{answer}}.",
		},
	}

	results := applyMustache(dep.Params, perform.Action, false)
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

	results := applyMustache(params, actionParams, true)

	assert.Equal(t,
		"To be, or not to be, that is the question.",
		strings.Join(results, " "))

}
