package docdb

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func (d *DocDB) newDoc() Doc {
	return Doc{
		filesPath: d.filesPath,
	}
}

// Path to docfile
func (d *Doc) filePath() string {
	return filepath.Join(d.dirPath(), d.Name)
}

// Path to directory where doc file is stored
func (d *Doc) dirPath() string {
	return filepath.Join(d.filesPath, strconv.Itoa(d.ID))
}

// Open Doc as a file
func (d *Doc) Open(doc *Doc) (*os.File, error) {
	return os.Open(d.filePath())
}

// Copies `file` to `<filesPath>/<id>/<filename>`
func (d *Doc) writeFile(file *os.File) error {
	err := os.Mkdir(d.dirPath(), 0750)
	if err != nil {
		return err
	}
	dest, err := os.Create(d.filePath())
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, file)
	return err
}

// Executes a `query` and parses the result as a slice of `Doc`s
func (d *DocDB) find(sqlStmt string, args ...any) ([]*Doc, error) {
	rows, err := d.db.Query(sqlStmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := []*Doc{}
	for rows.Next() {
		doc := d.newDoc()
		if err := rows.Scan(&doc.ID, &doc.Name, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}
