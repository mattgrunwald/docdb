package docdb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func getWd(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get working directory")
	}
	return wd
}
func getTestDbFile(t *testing.T, id string) string {
	return filepath.Join(getWd(t), fmt.Sprintf("%s_test.db", id))
}

func getTestDir(t *testing.T, id string) string {
	return filepath.Join(getWd(t), fmt.Sprintf("%s_test_db_files", id))
}

func openTestFile(t *testing.T, fileName string) *os.File {
	file, err := os.Open(filepath.Join(getWd(t), "test_files", fileName))
	if err != nil {
		t.Fatalf("failed to open test file %s\n", fileName)
	}
	return file
}

func setUpTest(t *testing.T) (*DocDB, func(t *testing.T), string) {
	id := uuid.New().String()
	dbFile := getTestDbFile(t, id)
	dir := getTestDir(t, id)
	db := new(dbFile, dir)
	tearDownTest := makeTearDownTest(id, db, dbFile, dir)
	return db, tearDownTest, id
}

func makeTearDownTest(id string, db *DocDB, dbFile, dir string) func(t *testing.T) {
	return func(t *testing.T) {
		err := db.db.Close()
		if err != nil {
			t.Logf("Failed to close DB: %q", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Logf("Failed to cleanup %s: %q", dir, err)
		}
		err = os.Remove(dbFile)
		if err != nil {
			t.Logf("Failed to cleanup %s: %q", dbFile, err)
		}
	}
}
