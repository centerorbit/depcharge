package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
)

func TestFindActionHandler(t *testing.T) {
	var handler func([]Dep, Perform)

	handler = findActionHandler("git")
	functionEqual(t, gitActionHandler, handler)
	handler = findActionHandler("secret")
	functionEqual(t, secretesActionHandler, handler)
	handler = findActionHandler("other")
	functionEqual(t, defaultActionHandler, handler)
}

func functionEqual(t *testing.T, expected func([]Dep, Perform), actual func([]Dep, Perform)){
	assert.Equal(t, reflect.ValueOf(expected), reflect.ValueOf(actual))
}