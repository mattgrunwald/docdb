package docdb_test

import (
	"os"
	"testing"

	"github.com/mattgrunwald/docdb"
)

func TestNew(t *testing.T) {
	db := docdb.New("testdb", "testfiles")

	t.Cleanup(func() {
		db.Close()
		os.Remove("testdb")
		os.RemoveAll("testfiles")
	})
}
