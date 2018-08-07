package main

import (
	"testing"
	"os"
	"flag"
)

func TestDepMain(t *testing.T) {
	oldArgs := os.Args

	os.Args = []string{"", "--kind=go", "--dryrun", "get", "{{get}}"}
	main()

	os.Args = oldArgs
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}
