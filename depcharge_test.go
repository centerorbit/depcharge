package main

import (
	"github.com/integrii/flaggy"
	"os"
	"testing"
)

func TestDepMain(t *testing.T) {
	oldArgs := os.Args

	os.Args = []string{"", "--kind=go", "--dryrun", "--", "get", "{{get}}"}
	main()

	os.Args = oldArgs
	flaggy.ResetParser()
}
