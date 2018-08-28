package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestExists(t *testing.T) {
	does, _ := exists("./docker")
	assert.True(t, does)

	does, _ = exists("./not-exists")
	assert.False(t, does)
}

func TestDepInjDefaultAction(t *testing.T) {
	result := depInjDefaultAction()
	assert.Equal(t, reflect.ValueOf(placeholderDefaultActionMock), reflect.ValueOf(result))

	complete := make(chan bool, 1)
	fmt.Println("The following errors can be ignored:")
	result(complete,Dep{},Perform{})
	fmt.Println("Okay, resume caring about errors.")
	<-complete
}

func TestAskForConfirmationYes(t *testing.T) {

	in, err := setupMockInput("Yes")
	if err != nil {
		t.Fatal(err)
	}
	result := askForConfirmation("Test Yes", in)
	assert.True(t, result)
	in.Close()
	fmt.Println()

	in, err = setupMockInput("y")
	if err != nil {
		t.Fatal(err)
	}
	result = askForConfirmation("Test y", in)
	assert.True(t, result)
	in.Close()
	fmt.Println()

	in, err = setupMockInput("YES")
	if err != nil {
		t.Fatal(err)
	}
	result = askForConfirmation("Test YES", in)
	assert.True(t, result)
	in.Close()
	fmt.Println()

	in, err = setupMockInput("yes")
	if err != nil {
		t.Fatal(err)
	}
	result = askForConfirmation("Test yes", in)
	assert.True(t, result)
	in.Close()
	fmt.Println()

}

func TestAskForConfirmationNo(t *testing.T) {

	in, err := setupMockInput("No")
	if err != nil {
		t.Fatal(err)
	}
	result := askForConfirmation("Test No", in)
	assert.False(t, result)
	in.Close()
	fmt.Println()

	in, err = setupMockInput("n")
	if err != nil {
		t.Fatal(err)
	}
	result = askForConfirmation("Test n", in)
	assert.False(t, result)
	in.Close()
	fmt.Println()


	in, err = setupMockInput("NO")
	if err != nil {
		t.Fatal(err)
	}
	result = askForConfirmation("Test NO", in)
	assert.False(t, result)
	in.Close()
	fmt.Println()


	in, err = setupMockInput("no")
	if err != nil {
		t.Fatal(err)
	}
	result = askForConfirmation("Test no", in)
	assert.False(t, result)
	in.Close()
	fmt.Println()

}


func TestAskForConfirmationAgain(t *testing.T) {
	in, err := setupMockInput("Nope \n n")
	if err != nil {
		t.Fatal(err)
	}
	result := askForConfirmation("", in)
	assert.False(t, result)
	in.Close()
	fmt.Println()
}


func TestPosString(t *testing.T) {
	haystack := []string{"Hello", "World", "Needle"}
	pos := posString(haystack, "Needle")
	assert.Equal(t, 2, pos)

	pos = posString(haystack, "Goose")
	assert.Equal(t, -1, pos)
}

func TestContainsString(t *testing.T) {
	haystack := []string{"Hello", "World", "Needle"}
	found := containsString(haystack, "Needle")
	assert.True(t, found)

	found = containsString(haystack, "Goose")
	assert.False(t, found)
}

func setupMockInput(response string) (in *os.File, err error)  {
	in, err = ioutil.TempFile("", "")
	if err != nil {
		return
	}

	_, err = io.WriteString(in, response)
	if err != nil {
		return
	}

	_, err = in.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	return
}