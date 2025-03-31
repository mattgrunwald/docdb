package docdb_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/mattgrunwald/docdb"
	"github.com/mattgrunwald/docdb/col"
	"github.com/mattgrunwald/docdb/order"
	"github.com/stretchr/testify/assert"
)

func openTestFile(t *testing.T, fileName string) *os.File {
	t.Helper()
	file, err := os.Open(filepath.Clean(filepath.Join("test_files", fileName)))
	if err != nil {
		t.Fatalf("failed to open test file %s\n", fileName)
	}
	return file
}

func setUpTest(t *testing.T) (*docdb.DocDB, string) {
	t.Helper()
	dbDir := t.TempDir()
	dbFile, _ := os.CreateTemp(dbDir, "")
	defer dbFile.Close()
	fileDir := t.TempDir()
	db := docdb.New(dbFile.Name(), fileDir)
	t.Cleanup(func() {
		err := db.Close()
		if err != nil {
			t.Logf("Failed to close DB: %q", err)
		}
	})
	return db, fileDir
}

func compareTimes(t *testing.T, expected time.Time, received time.Time) {
	assert.Equal(t, expected.Year(), received.Year())
	assert.Equal(t, expected.Month(), received.Month())
	assert.Equal(t, expected.Day(), received.Day())
	assert.Equal(t, expected.Hour(), received.Hour())
	assert.Equal(t, expected.Minute(), received.Minute())
	assert.Equal(t, expected.Second(), received.Second())
}

func compareDocs(t *testing.T, expected *docdb.Doc, received *docdb.Doc) {
	assert.Equal(t, expected.ID, received.ID)
	assert.Equal(t, expected.Name, received.Name)
	compareTimes(t, expected.CreatedAt.UTC(), received.CreatedAt.UTC())
	compareTimes(t, expected.UpdatedAt.UTC(), received.UpdatedAt.UTC())
}

func Test_Insert(t *testing.T) {
	// setup
	db, _ := setUpTest(t)

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
	db, filesDir := setUpTest(t)
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
	assert.Equal(t, 1, doc.ID)
	assert.Equal(t, "b.txt", doc.Name)
	compareTimes(t, insertTime.UTC(), doc.CreatedAt)
	compareTimes(t, updateTime.UTC(), doc.UpdatedAt)

	files, _ := os.ReadDir(filepath.Join(filesDir, strconv.Itoa(doc.ID)))
	assert.Equal(t, 1, len(files))
	assert.Equal(t, "b.txt", filepath.Base(files[0].Name()))
}

func Test_findOne(t *testing.T) {
	t.Run("Empty slice when there are no results", func(t *testing.T) {
		// setup
		db, _ := setUpTest(t)

		// test case
		doc, err := db.FindOne(0)

		// checks
		assert.Nil(t, doc)
		assert.NotNil(t, err)
	})

	t.Run("Result is complete correct", func(t *testing.T) {
		// setup
		db, _ := setUpTest(t)
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
		db, _ := setUpTest(t)

		docs, err := db.FindMany(2, 0, col.UpdatedAt, order.ASC)
		if err != nil {
			t.Logf("findMany failed: %q\n", err)
			t.Fail()
		}
		assert.Equal(t, 0, len(docs))
	})

	t.Run("results are complete and in correct order", func(t *testing.T) {
		// setup
		db, _ := setUpTest(t)
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
		db, _ := setUpTest(t)

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
		db, _ := setUpTest(t)
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
				t.Logf("FindAll failed: %q\n", err)
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
				t.Logf("FindAll failed: %q\n", err)
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

func Test_FindLike(t *testing.T) {
	// setup
	db, _ := setUpTest(t)
	fileA := openTestFile(t, "a.txt")
	fileB := openTestFile(t, "b.txt")
	fileC := openTestFile(t, "c.txt")
	defer fileA.Close()
	defer fileB.Close()
	defer fileC.Close()
	docA, _ := db.Insert(fileA)
	docB, _ := db.Insert(fileB)
	docC, _ := db.Insert(fileC)

	t.Run("a", func(t *testing.T) {
		docs, err := db.FindLike("a")
		if err != nil {
			t.Logf("FindLike failed: %q\n", err)
			t.Fail()
		}

		assert.Equal(t, 1, len(docs))
		assert.Equal(t, docA.Name, docs[0].Name)
	})

	t.Run(".txt", func(t *testing.T) {
		docs, err := db.FindLike(".txt")
		if err != nil {
			t.Logf("FindLike failed: %q\n", err)
			t.Fail()
		}

		assert.Equal(t, 3, len(docs))
		assert.Equal(t, docA.Name, docs[0].Name)
		assert.Equal(t, docB.Name, docs[1].Name)
		assert.Equal(t, docC.Name, docs[2].Name)
	})
}

func Test_Delete(t *testing.T) {
	// setup
	db, _ := setUpTest(t)
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
