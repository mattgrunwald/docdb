package docdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mattgrunwald/docdb/col"
	"github.com/mattgrunwald/docdb/order"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.RemoveAll("test_db_files")
	os.Remove("test.db")
}

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
	file, err := os.Open(filepath.Clean(filepath.Join(getWd(t), "test_files", fileName)))
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

func compareTimes(t *testing.T, expected time.Time, received time.Time) {
	assert.Equal(t, expected.Year(), received.Year())
	assert.Equal(t, expected.Month(), received.Month())
	assert.Equal(t, expected.Day(), received.Day())
	assert.Equal(t, expected.Hour(), received.Hour())
	assert.Equal(t, expected.Minute(), received.Minute())
	assert.Equal(t, expected.Second(), received.Second())
}

func compareDocs(t *testing.T, expected *Doc, received *Doc) {
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Name, received.Name)
	compareTimes(t, expected.CreatedAt.UTC(), received.CreatedAt.UTC())
	compareTimes(t, expected.UpdatedAt.UTC(), received.UpdatedAt.UTC())
}

func Test_Insert(t *testing.T) {
	// setup
	db, tearDownTest, _ := setUpTest(t)
	defer tearDownTest(t)

	file := openTestFile(t, "a.txt")
	defer file.Close()

	// test case
	doc, err := db.Insert(file)
	if err != nil {
		t.Logf("Insertion failed: %q\n", err)
		t.Fail()
	}

	// checks
	assert.Equal(t, 1, doc.ID)
	assert.Equal(t, "a.txt", doc.Name)
	compareTimes(t, time.Now().UTC(), doc.CreatedAt)
	compareTimes(t, time.Now().UTC(), doc.UpdatedAt)
}

func Test_Update(t *testing.T) {
	// setup
	db, tearDownTest, testID := setUpTest(t)
	defer tearDownTest(t)
	fileA := openTestFile(t, "a.txt")
	defer fileA.Close()
	fileB := openTestFile(t, "b.txt")
	defer fileB.Close()

	insertTime := time.Now()
	doc, err := db.Insert(fileA)
	if err != nil {
		t.Logf("Insertion failed: %q\n", err)
		t.Fail()
	}

	// test case
	doc, err = db.Update(doc.ID, fileB)
	if err != nil {
		t.Logf("Update failed: %q\n", err)
		t.Fail()
	}
	updateTime := time.Now()

	// checks
	rows, _ := db.db.Query("SELECT id, name, created_at, updated_at FROM docs")
	defer rows.Close()
	rows.Next()
	err = rows.Scan(&doc.ID, &doc.Name, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		t.Logf("Querying DB failed: %q\n", err)
		t.Fail()
	}
	assert.Equal(t, 1, doc.ID)
	assert.Equal(t, "b.txt", doc.Name)
	compareTimes(t, insertTime.UTC(), doc.CreatedAt)
	compareTimes(t, updateTime.UTC(), doc.UpdatedAt)

	files, _ := os.ReadDir(filepath.Join(getTestDir(t, testID), strconv.Itoa(doc.ID)))
	assert.Equal(t, 1, len(files))
	assert.Equal(t, "b.txt", filepath.Base(files[0].Name()))
}

func Test_findOne(t *testing.T) {
	t.Run("Empty slice when there are no results", func(t *testing.T) {
		// setup
		db, tearDownTest, _ := setUpTest(t)
		defer tearDownTest(t)

		// test case
		doc, err := db.FindOne(0)

		// checks
		assert.Nil(t, doc)
		assert.NotNil(t, err)
	})

	t.Run("Result is complete correct", func(t *testing.T) {
		// setup
		db, tearDownTest, _ := setUpTest(t)
		defer tearDownTest(t)
		fileA := openTestFile(t, "a.txt")
		defer fileA.Close()
		_, _ = db.Insert(fileA)
		fileB := openTestFile(t, "b.txt")
		defer fileB.Close()
		docB, _ := db.Insert(fileB)

		// test case
		doc, err := db.FindOne(docB.ID)
		if err != nil {
			t.Logf("findOne failed: %q\n", err)
			t.Fail()
		}

		// checks
		compareDocs(t, docB, doc)
	})
}

func Test_FindMany(t *testing.T) {
	t.Run("Empty slice when there are no results", func(t *testing.T) {
		// setup
		db, tearDownTest, _ := setUpTest(t)
		defer tearDownTest(t)

		docs, err := db.FindMany(2, 0, col.UpdatedAt, order.ASC)
		if err != nil {
			t.Logf("findMany failed: %q\n", err)
			t.Fail()
		}
		assert.Equal(t, 0, len(docs))
	})

	t.Run("results are complete and in correct order", func(t *testing.T) {
		// setup
		db, tearDownTest, _ := setUpTest(t)
		defer tearDownTest(t)
		fileA := openTestFile(t, "a.txt")
		fileB := openTestFile(t, "b.txt")
		fileC := openTestFile(t, "c.txt")
		defer fileA.Close()
		defer fileB.Close()
		defer fileC.Close()
		docA, _ := db.Insert(fileA)
		docB, _ := db.Insert(fileB)
		_, _ = db.Insert(fileC)

		t.Run("ASC", func(t *testing.T) {
			// test case
			docs, err := db.FindMany(2, 0, col.ID, order.ASC)
			if err != nil {
				t.Logf("findMany failed: %q\n", err)
				t.Fail()
			}

			// checks
			assert.Equal(t, 2, len(docs))
			compareDocs(t, docA, docs[0])
			compareDocs(t, docB, docs[1])
		})
		t.Run("DESC", func(t *testing.T) {
			// test case
			docs, err := db.FindMany(2, 0, col.ID, order.DESC)
			if err != nil {
				t.Logf("findMany failed: %q\n", err)
				t.Fail()
			}

			// checks
			assert.Equal(t, 2, len(docs))
			compareDocs(t, docB, docs[1])
			compareDocs(t, docA, docs[0])
		})
	})
}

func Test_FindAll(t *testing.T) {
	t.Run("Empty slice when there are no results", func(t *testing.T) {
		// setup
		db, tearDownTest, _ := setUpTest(t)
		defer tearDownTest(t)

		//test case
		docs, err := db.FindAll(col.CreatedAt, order.ASC)
		if err != nil {
			t.Logf("findMany failed: %q\n", err)
			t.Fail()
		}

		// checks
		assert.Equal(t, 0, len(docs))
	})

	t.Run("results are complete and in correct order ASC", func(t *testing.T) {
		// setup
		db, tearDownTest, _ := setUpTest(t)
		defer tearDownTest(t)
		fileA := openTestFile(t, "a.txt")
		fileB := openTestFile(t, "b.txt")
		fileC := openTestFile(t, "c.txt")
		defer fileA.Close()
		defer fileB.Close()
		defer fileC.Close()
		docA, _ := db.Insert(fileA)
		docB, _ := db.Insert(fileB)
		docC, _ := db.Insert(fileC)

		t.Run("ASC", func(t *testing.T) {
			// test case
			docs, err := db.FindAll(col.Name, order.ASC)
			if err != nil {
				t.Logf("findMany failed: %q\n", err)
				t.Fail()
			}

			// checks
			assert.Equal(t, len(docs), 3)
			compareDocs(t, docA, docs[0])
			compareDocs(t, docB, docs[1])
			compareDocs(t, docC, docs[2])
		})

		t.Run("DESC", func(t *testing.T) {
			// test case
			docs, err := db.FindAll(col.Name, order.DESC)
			if err != nil {
				t.Logf("findMany failed: %q\n", err)
				t.Fail()
			}

			// checks
			assert.Equal(t, len(docs), 3)
			compareDocs(t, docC, docs[2])
			compareDocs(t, docB, docs[1])
			compareDocs(t, docA, docs[0])
		})

	})
}

func Test_Delete(t *testing.T) {
	// setup
	db, tearDownTest, _ := setUpTest(t)
	defer tearDownTest(t)
	fileA := openTestFile(t, "a.txt")
	fileB := openTestFile(t, "b.txt")
	defer fileA.Close()
	defer fileB.Close()
	docA, _ := db.Insert(fileA)
	docB, _ := db.Insert(fileB)

	// test case
	err := db.Delete(docA.ID)
	if err != nil {
		t.Logf("Delete failed: %q\n", err)
		t.Fail()
	}

	//checks
	docs, _ := db.FindAll(col.Name, order.ASC)
	assert.Equal(t, 1, len(docs))
	compareDocs(t, docB, docs[0])

}
