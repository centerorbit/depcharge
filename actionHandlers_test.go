package main

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFindActionHandler(t *testing.T) {
	var handler func(chan bool, []Dep, Perform) int

	handler = findActionHandler("git")
	assert.Equal(t, reflect.ValueOf(gitActionHandler), reflect.ValueOf(handler))
	handler = findActionHandler("secret")
	assert.Equal(t, reflect.ValueOf(secretesActionHandler), reflect.ValueOf(handler))
	handler = findActionHandler("other")
	assert.Equal(t, reflect.ValueOf(defaultActionHandler), reflect.ValueOf(handler))
}

func TestGitHandlerClone(t *testing.T) {
	defer func() {
		mockDefaultAction = nil
	}()

	called := 0
	mockDefaultAction = func(complete chan<- bool, dep Dep, perform Perform) {
		assert.Equal(t, "git", dep.Kind)
		assert.Equal(t, "clone", perform.Action[0])
		assert.Equal(t, "source", perform.Action[1])
		assert.Equal(t, "location", perform.Action[2])
		assert.True(t, perform.DryRun)
		called++
		complete <- true
	}

	perform := Perform{
		Action: []string{"clone", "source", "location"},
		DryRun: true,
	}

	deps := []Dep{
		{
			Kind: "git",
		},
	}

	complete := make(chan bool)
	n := gitActionHandler(complete, deps, perform)

	drainChannel(n, complete)

	assert.Equal(t, 1, n)
	assert.Equal(t, 1, called)
}

func TestGitHandlerStatus(t *testing.T) {
	defer func() {
		mockDefaultAction = nil
	}()

	called := 0
	mockDefaultAction = func(complete chan<- bool, dep Dep, perform Perform) {
		assert.Equal(t, "git", dep.Kind)
		assert.Equal(t, "status", perform.Action[0])
		assert.False(t, perform.DryRun)
		called++
		complete <- true
	}

	perform := Perform{
		Action: []string{"status"},
	}

	deps := []Dep{
		{
			Kind: "git",
		},
	}

	complete := make(chan bool)
	n := gitActionHandler(complete, deps, perform)

	drainChannel(n, complete)

	assert.Equal(t, 1, n)
	assert.Equal(t, 1, called)
}

func TestSecretsHandler(t *testing.T) {
	defer func() { mockDefaultAction = nil }()

	called := 0
	mockDefaultAction = func(complete chan<- bool, dep Dep, perform Perform) {
		assert.Equal(t, "secret", dep.Kind)
		assert.Equal(t, "doesn't", perform.Action[0])
		assert.Equal(t, "matter", perform.Action[1])
		called++
		complete <- true
	}

	perform := Perform{
		Action: []string{"doesn't", "matter", "...yet"},
		DryRun: true,
	}

	deps := []Dep{
		{
			Kind: "secret",
		},
		{
			Kind: "secret",
		},
	}

	complete := make(chan bool)
	n := secretesActionHandler(complete, deps, perform)

	drainChannel(n, complete)

	assert.Equal(t, 2, n)
	assert.Equal(t, 2, called)
}
