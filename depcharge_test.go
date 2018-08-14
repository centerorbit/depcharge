package main

import (
	"testing"
	"os"
	"github.com/integrii/flaggy"
)

func TestDepMain(t *testing.T) {
	oldArgs := os.Args

	os.Args = []string{"", "--kind=go", "--dryrun", "--", "get", "{{get}}"}
	main()

	os.Args = oldArgs
	flaggy.ResetParser()
}
