package docdb

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func (d *DocDB) filePath(id int, name string) string {
	return filepath.Join(d.dirPath(id), name)
}

func (d *DocDB) dirPath(id int) string {
	return filepath.Join(d.filesPath, strconv.Itoa(id))
}

func (d *DocDB) docToFile(doc *Doc) (*os.File, error) {
	return os.Open(d.filePath(doc.ID, doc.Name))
}

func (d *DocDB) docsToFiles(docs []*Doc) ([]*os.File, error) {
	files := []*os.File{}
	for _, doc := range docs {
		file, err := d.docToFile(doc)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (d *DocDB) writeFile(id int, file *os.File) error {
	name := filepath.Base(file.Name())
	err := os.Mkdir(d.dirPath(id), os.ModePerm)
	if err != nil {
		return err
	}
	dest, err := os.Create(d.filePath(id, name))
	if err != nil {
		return err
	}
	_, err = io.Copy(dest, file)
	return err
}

func (d *DocDB) find(sqlStmt string) ([]*Doc, error) {
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docs := []*Doc{}
	for rows.Next() {
		var doc Doc
		if err := rows.Scan(&doc.ID, &doc.Name, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}
