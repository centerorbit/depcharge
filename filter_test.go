package main

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestUnwrap(t *testing.T) {

	var deps []dep
	var labels []string
	var foundDeps []dep

	// Test if nothing is okay
	foundDeps = unwrap(deps, "", labels, nil)
	assert.Empty(t, foundDeps)

	testDep := dep{
		Name:     "sample",
		Kind:     "git",
		Location: "sample",
		DepList:  nil,
		Labels:   nil,
		Params:   nil,
	}

	// Tests single-level dep
	deps = append(deps, testDep)
	deps[0].Params = map[string]string{"kind": testDep.Kind, "location": testDep.Location, "name": testDep.Name}
	foundDeps = unwrap(deps, "./", labels, nil)
	assert.Equal(t, deps, foundDeps)

	//Making a nest
	testDep = dep{
		Name:     "parent",
		Kind:     "git",
		Location: "parent-dir",
		DepList: []dep{
			{
				Name:     "child",
				Kind:     "git",
				Location: "child-dir",
				Labels:   []string{"first"},
				DepList: []dep{
					{
						Name:     "grandchild",
						Kind:     "git",
						Location: "grandchild-dir",
						Labels:   []string{"second"},
					},
				},
			},
			{
				Name:     "sibling",
				Kind:     "git",
				Location: "sample",
			},
		},
	}

	deps = append(deps, testDep)

	foundDeps = unwrap(deps, "./append/", labels, nil)

	// Test flattening of deps:
	assert.Equal(t, 5, len(foundDeps))

	// Test dir expanding
	assert.Equal(t, filepath.Clean("./append/parent-dir"), foundDeps[4].Location)                          // First in, least nested, last in array
	assert.Equal(t, filepath.Clean("./append/sample"), foundDeps[0].Location)                              // Last in, most nested
	assert.Equal(t, filepath.Clean("./append/parent-dir/child-dir/grandchild-dir"), foundDeps[1].Location) // Last in, most nested

	// Test label inheritance
	assert.Empty(t, foundDeps[4].Labels)
	assert.Equal(t, []string{"first"}, foundDeps[2].Labels)
	assert.Equal(t, []string{"second", "first"}, foundDeps[1].Labels)
}

func TestApplyFilterKind(t *testing.T) {
	testDeps := []dep{
		{
			Labels: nil,
		},
		{
			Labels: []string{"one"},
		},
		{
			Labels: []string{"one", "two"},
		},
		{
			Labels: []string{"three", "two"},
		},
	}

	perform := perform{
		Labels:    "",
		Exclusive: false,
		Force:     true,
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
	testDeps := []dep{
		{
			Kind: "git",
		},
		{
			Kind: "git",
		},
		{
			Kind: "git",
		},
		{
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

func TestIsExclusive(t *testing.T) {
	what := []string{"one", "two"} // dep.labels
	against := []string{"one"}     // labels

	// dep does have "one"
	assert.True(t, isExclusive(what, against))

	// dep does not have "three"
	against = []string{"one", "three"} // labels
	assert.False(t, isExclusive(what, against))

	// Testing no labels, no filter, all comes back
	against = []string{} // labels
	assert.True(t, isExclusive(what, against))

	// Testing no dep.labels, and no labels
	what = []string{} // dep.labels
	assert.True(t, isExclusive(what, against))

	// Testing no dep.labels, but with labels
	against = []string{"one", "three"} // labels
	assert.False(t, isExclusive(what, against))
}

func TestIsInclusive(t *testing.T) {
	what := []string{"one", "two"} // dep.labels
	against := []string{"one"}     // labels

	// dep does have "one"
	assert.True(t, isInclusive(what, against))

	// dep does not have "three", but does have "one"
	against = []string{"one", "three"} // labels
	assert.True(t, isInclusive(what, against))

	// dep does not have "four"
	against = []string{"four"} // labels
	assert.False(t, isInclusive(what, against))

	// Testing no labels, no filter, all comes back
	against = []string{} // labels
	assert.True(t, isInclusive(what, against))

	// Testing no dep.labels, and no labels
	what = []string{} // dep.labels
	assert.True(t, isInclusive(what, against))

	// Testing no dep.labels, but with labels
	against = []string{"one", "three"} // labels
	assert.False(t, isInclusive(what, against))
}
