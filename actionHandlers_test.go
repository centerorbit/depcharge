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

func functionEqual(t *testing.T,
	expected func(chan<- bool, []Dep, Perform) int,
	actual func(chan<- bool, []Dep, Perform) int) {
	assert.Equal(t, reflect.ValueOf(expected), reflect.ValueOf(actual))
}
