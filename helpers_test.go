package docdb_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Open(t *testing.T) {
	// setup
	db, _ := setUpTest(t)

	inputFile := openTestFile(t, "a.txt")
	defer inputFile.Close()
	doc, _ := db.Insert(inputFile)

	// test case
	outputFile, err := doc.Open()
	if err != nil {
		t.Logf("Open file: %q\n", err)
		t.Fail()
	}
	defer outputFile.Close()

	// checks
	assert.Equal(t, doc.Name, filepath.Base(outputFile.Name()))
}
