package docdb_test

import (
	"os"
	"testing"

	"github.com/mattgrunwald/docdb"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	filesPath := "testfiles"
	dbFile := "testdb"
	_, err := os.Stat(filesPath)
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(dbFile)
	assert.True(t, os.IsNotExist(err))
	db := docdb.New(dbFile, filesPath)

	t.Cleanup(func() {
		db.Close()
		os.Remove("testdb")
		os.RemoveAll("testfiles")
	})
}
