package main

import (
	"github.com/integrii/flaggy"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

/**
This will load the dep.yml file of depcharge
*/
func TestLoad(t *testing.T) {
	results := load()
	assert.NotEmpty(t, results)
	assert.Equal(t, 1, len(results.Deps))
	assert.True(t, len(results.Deps[0].DepList) > 3)
}

func TestProcessArgs(t *testing.T) {
	oldArgs := os.Args

	os.Args = []string{"", "--kind=git", "--inclusive", "--dryrun", "--labels=some,thing", "--", "status"}
	results := processArgs()

	assert.Equal(t, "git", results.Kind)
	assert.Equal(t, "some,thing", results.Labels)
	assert.Equal(t, "status", results.Action[0])
	assert.True(t, results.DryRun)
	assert.False(t, results.Exclusive)

	os.Args = oldArgs
	flaggy.ResetParser()
}
