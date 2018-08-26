package main

import (
	"github.com/integrii/flaggy"
	"os"
	"testing"
)

const COVER_LIMIT = 0.8

//func TestMain(m *testing.M) {
//	// call flag.Parse() here if TestMain uses flags
//	rc := m.Run()
//
//	// rc 0 means we've passed,
//	// and CoverMode will be non empty if run with -cover
//	if rc == 0 && testing.CoverMode() != "" {
//		c := testing.Coverage()
//		if c < COVER_LIMIT && os.LookupEnv("COVER_STRICT") == "true" {
//			fmt.Println("Tests passed but coverage was below ",COVER_LIMIT*100,"%")
//			rc = -1
//		}
//	}
//	os.Exit(rc)
//}

func TestDepMain(t *testing.T) {
	oldArgs := os.Args

	os.Args = []string{"", "--kind=go", "--force", "--dryrun", "--", "get", "{{get}}"}
	main()

	flaggy.ResetParser()

	os.Args = []string{"", "--force", "--", "go", "get", "{{get}}"}
	main()

	os.Args = oldArgs
	flaggy.ResetParser()
}
