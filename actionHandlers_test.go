package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFindActionHandler(t *testing.T) {
	var handler func(chan<- bool, []Dep, Perform) int

	handler = findActionHandler("git")
	functionEqual(t, gitActionHandler, handler)
	handler = findActionHandler("secret")
	functionEqual(t, secretesActionHandler, handler)
	handler = findActionHandler("other")
	functionEqual(t, defaultActionHandler, handler)
}

func TestGitHandler(t *testing.T) {
	perform := Perform{
		Action: []string{"clone", "source", "location"},
		DryRun: true,
	}

	deps := []Dep{
		{},
	}

	complete := make(chan bool)
	n := gitActionHandler(complete, deps, perform)
	assert.Equal(t, 1, n)

	perform = Perform{
		Action: []string{"status"},
		DryRun: true,
	}
	n = gitActionHandler(complete, deps, perform)
	assert.Equal(t, 1, n)
}

func TestSecretsHandler(t *testing.T) {
	perform := Perform{
		Action: []string{"clone", "source", "location"},
		DryRun: true,
	}

	deps := []Dep{
		{},
	}

	complete := make(chan bool)
	n := secretesActionHandler(complete, deps, perform)
	assert.Equal(t, 1, n)
}

func functionEqual(t *testing.T,
	expected func(chan<- bool, []Dep, Perform) int,
	actual func(chan<- bool, []Dep, Perform) int) {
	assert.Equal(t, reflect.ValueOf(expected), reflect.ValueOf(actual))
}
