package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTemplateParams(t *testing.T) {

	dep := Dep{
		Kind:     "git",
		Name:     "depcharge",
		Location: "./",
		Params: map[string]string{
			"is":     "be",
			"is not": "not to be",
			"answer": "question",
		},
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
