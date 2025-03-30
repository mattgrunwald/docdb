package docdb

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattgrunwald/docdb/col"
	"github.com/mattgrunwald/docdb/order"
)

// Insert a file
func (d *DocDB) Insert(file *os.File) (*Doc, error) {
	doc := d.newDoc()
	doc.Name = filepath.Base(file.Name())

	sqlStmt := fmt.Sprintf(
		`INSERT INTO 
			docs
			(name, created_at, updated_at)
		VALUES
			('%s', datetime('now'), datetime('now'));`,
		doc.Name)
	res, err := d.db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	insertedDoc, err := d.FindOne(int(id))
	if err != nil {
		return nil, err
	}
	err = insertedDoc.writeFile(file)
	if err != nil {
		return nil, err
	}
	return insertedDoc, nil
}

// Update a file
func (d *DocDB) Update(id int, file *os.File) (*Doc, error) {
	fileName := filepath.Base(file.Name())
	sqlStmt := fmt.Sprintf(`
		UPDATE docs 
		SET 
			name = '%s',
			updated_at = datetime('now')
		WHERE id = '%d';`,
		fileName, id)
	_, err := d.db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	doc, err := d.FindOne(id)
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(doc.dirPath())
	if err != nil {
		return nil, err
	}
	err = doc.writeFile(file)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Find on file by its ID
func (d *DocDB) FindOne(id int) (*Doc, error) {
	sqlStmt := fmt.Sprintf(`
		SELECT id, name, created_at, updated_at 
		FROM docs 
		WHERE id = %d 
		LIMIT 1;`,
		id,
	)
	rows, err := d.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	found := rows.Next()
	if !found {
		return nil, fmt.Errorf("Row with id %d not found", id)
	}
	doc := d.newDoc()
	if err = rows.Scan(&doc.ID, &doc.Name, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
		return nil, err
	}
	return &doc, nil
}

// Find many files. Useful for pagination.
func (d *DocDB) FindMany(count int, offset int, orderCol col.DocCol, order order.Order) ([]*Doc, error) {
	sqlStmt := fmt.Sprintf(`
		SELECT 
			id, name, created_at, updated_at
		FROM docs 
		ORDER BY %s %s
		LIMIT %d
		OFFSET %d;`,
		orderCol, order, count, offset,
	)
	return d.find(sqlStmt)
}

// Find all files.
func (d *DocDB) FindAll(orderCol col.DocCol, order order.Order) ([]*Doc, error) {
	sqlStmt := fmt.Sprintf(`
		SELECT 
			id, name, created_at, updated_at 
		FROM docs
		ORDER BY %s %s;`,
		orderCol, order,
	)
	return d.find(sqlStmt)
}

func (d *DocDB) Delete(id int) error {
	sqlStmt := fmt.Sprintf(`
		DELETE 
		FROM docs 
		WHERE id = %d;`,
		id,
	)
	_, err := d.db.Exec(sqlStmt)
	return err
}
