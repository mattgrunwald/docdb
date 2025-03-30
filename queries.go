package docdb

import (
	"os"
	"path/filepath"

	"github.com/mattgrunwald/docdb/col"
	"github.com/mattgrunwald/docdb/order"
)

// Insert a file
func (d *DocDB) Insert(file *os.File) (*Doc, error) {
	doc := d.newDoc()
	doc.Name = filepath.Base(file.Name())

	res, err := d.db.Exec(`
		INSERT INTO 
			docs
			(name, created_at, updated_at)
		VALUES
			(?, datetime('now'), datetime('now'));`,
		doc.Name,
	)
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
	_, err := d.db.Exec(`
		UPDATE docs 
		SET 
			name = ?,
			updated_at = datetime('now')
		WHERE id = ?;`,
		fileName, id,
	)
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
	doc := d.newDoc()
	err := d.db.QueryRow(`
		SELECT id, name, created_at, updated_at 
		FROM docs 
		WHERE id = ?;`,
		id,
	).Scan(&doc.ID, &doc.Name, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// Find many files. Useful for pagination.
func (d *DocDB) FindMany(count int, offset int, orderCol col.DocCol, ord order.Order) ([]*Doc, error) {
	if ord == order.ASC {
		return d.find(`
		SELECT 
			id, name, created_at, updated_at
		FROM docs 
		ORDER BY $1 ASC
		LIMIT $2
		OFFSET $3;`,
			orderCol, count, offset)
	}

	return d.find(`
		SELECT 
			id, name, created_at, updated_at
		FROM docs 
		ORDER BY $1 DESC
		LIMIT $2
		OFFSET $3;`,
		orderCol, count, offset,
	)

}

// Find all files.
func (d *DocDB) FindAll(orderCol col.DocCol, ord order.Order) ([]*Doc, error) {

	if ord == order.ASC {
		return d.find(`
			SELECT 
				id, name, created_at, updated_at 
			FROM docs
			ORDER BY $1 ASC;`,
			orderCol,
		)
	}
	return d.find(`
		SELECT 
			id, name, created_at, updated_at 
		FROM docs
		ORDER BY $1 DESC;`,
		orderCol)
}

func (d *DocDB) Delete(id int) error {
	_, err := d.db.Exec(`
		DELETE 
		FROM docs 
		WHERE id = ?;`,
		id,
	)
	return err
}
