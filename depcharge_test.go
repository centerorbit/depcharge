package main

import (
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const CoverLimit = 0.8

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		val, strict := os.LookupEnv("COVER_STRICT")

		if val == "false" || val == "0" {
			strict = false
		}

		if strict {
			if c < CoverLimit {
				fmt.Println("Tests passed but coverage was below ", CoverLimit*100, "%")
				rc = -1
			} else {
				fmt.Println("Strict coverage passed!")
			}
		}
	}
	os.Exit(rc)
}

var oldArgs []string

func TestDepMainDryRun(t *testing.T) {
	oldArgs = os.Args
	defer flaggy.ResetParser()
	defer func() { os.Args = oldArgs }()
	defer func() { mockDefaultAction = nil }()

	called := 0
	mockDefaultAction = func(complete chan<- bool, dep dep, perform perform) {
		assert.Equal(t, "go", perform.Kind)
		assert.True(t, perform.Force)
		assert.True(t, perform.DryRun)
		assert.Equal(t, "get", perform.Action[0])
		assert.Equal(t, "{{get}}", perform.Action[1])
		called++
		complete <- true
	}

	os.Args = []string{"", "--kind=go", "--force", "--dryrun", "--", "get", "{{get}}"}
	main()

	assert.Equal(t, 3, called)
}

func TestDepMainForce(t *testing.T) {
	oldArgs = os.Args
	defer flaggy.ResetParser()
	defer func() { os.Args = oldArgs }()
	defer func() { mockDefaultAction = nil }()

	called := 0
	mockDefaultAction = func(complete chan<- bool, dep dep, perform perform) {
		assert.Equal(t, "go", perform.Kind)
		assert.True(t, perform.Force)
		assert.False(t, perform.DryRun)
		assert.Equal(t, "get", perform.Action[0])
		assert.Equal(t, "{{get}}", perform.Action[1])
		called++
		complete <- true
	}

	os.Args = []string{"", "--force", "--", "go", "get", "{{get}}"}
	main()

	assert.Equal(t, 3, called)
}
