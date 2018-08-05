package main

import (
	"testing"
	"depcharge/src/github.com/stretchr/testify/assert"
		"path/filepath"
	)

func TestUnwrap(t *testing.T) {

	var deps []Dep
	var labels []string
	var foundDeps []Dep

	// Test if nothing is okay
	foundDeps = unwrap(deps,"", labels)
	assert.Empty(t, foundDeps)

	testDep := Dep{
		Name: "sample",
		Kind: "git",
		Location: "sample",
		DepList: nil,
		Labels: nil,
		Params: nil,
	}

	// Tests single-level dep
	deps = append(deps, testDep)
	foundDeps = unwrap(deps,"./", labels)
	assert.Equal(t, deps, foundDeps)

	//Making a nest
	testDep = Dep{
		Name: "parent",
		Kind: "git",
		Location: "parent-dir",
		DepList: []Dep{
			Dep{
				Name:     "child",
				Kind:     "git",
				Location: "child-dir",
				Labels:   []string{"first"},
				DepList:  []Dep{
					Dep{
						Name:     "grandchild",
						Kind:     "git",
						Location: "grandchild-dir",
						Labels:   []string{"second"},
					},
				},
			},
			Dep{
				Name:     "sibling",
				Kind:     "git",
				Location: "sample",
			},
		},
	}

	deps = append(deps, testDep)

	foundDeps = unwrap(deps,"./append/", labels)

	// Test flattening of deps:
	assert.Equal(t, 5, len(foundDeps))

	// Test dir expanding
	assert.Equal(t, filepath.Clean("./append/parent-dir"), foundDeps[4].Location) // First in, least nested, last in array
	assert.Equal(t, filepath.Clean("./append/sample"), foundDeps[0].Location)// Last in, most nested
	assert.Equal(t, filepath.Clean("./append/parent-dir/child-dir/grandchild-dir"), foundDeps[1].Location)// Last in, most nested


	// Test label inheritance
	assert.Empty(t, foundDeps[4].Labels)
	assert.Equal(t,[]string{"first"}, foundDeps[2].Labels)
	assert.Equal(t,[]string{"second", "first"}, foundDeps[1].Labels)
}

func TestApplyFilterKind(t *testing.T) {
	testDeps := []Dep{
		Dep{
			Labels: nil,
		},
		Dep{
			Labels: []string{"one"},
		},
		Dep{
			Labels: []string{"one", "two"},
		},
		Dep{
			Labels: []string{"three", "two"},
		},
	}

	perform := Perform{
		Labels: "",
		Exclusive: false,
		Force: true,
	}

	result := applyFilterLabel(testDeps, perform)
	assert.Equal(t, 4, len(result))

	perform.Labels = "one"
	perform.Exclusive = true
	result = applyFilterLabel(testDeps, perform)
	assert.Equal(t, 2, len(result))

	perform.Exclusive = false
	result = applyFilterLabel(testDeps, perform)
	assert.Equal(t, 2, len(result))

	perform.Labels = "one,two"
	perform.Exclusive = true
	result = applyFilterLabel(testDeps, perform)
	assert.Equal(t, 1, len(result))

	perform.Labels = "two,one"
	perform.Exclusive = false
	result = applyFilterLabel(testDeps, perform)
	assert.Equal(t, 3, len(result))

	perform.Labels = "zero"
	perform.Exclusive = false
	result = applyFilterLabel(testDeps, perform)
	assert.Equal(t, 0, len(result))
}


func TestApplyFilterLabel(t *testing.T) {
	testDeps := []Dep{
		Dep{
			Kind: "git",
		},
		Dep{
			Kind: "git",
		},
		Dep{
			Kind: "git",
		},
		Dep{
			Kind: "other",
		},
	}

	result := applyFilterKind(testDeps, "git")
	assert.Equal(t, 3, len(result))

	result = applyFilterKind(testDeps, "other")
	assert.Equal(t, 1, len(result))

	result = applyFilterKind(testDeps, "none")
	assert.Equal(t, 0, len(result))
}