package docdb

import (
	"fmt"
	"os"
	"path/filepath"
)

func (d *DocDB) Insert(file *os.File) (*Doc, error) {
	doc := Doc{
		Name: filepath.Base(file.Name()),
	}
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
	err = d.writeFile(int(id), file)
	if err != nil {
		return nil, err
	}
	return d.findOne(int(id))
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
	err = os.RemoveAll(d.dirPath(id))
	if err != nil {
		return nil, err
	}
	err = d.writeFile(int(id), file)
	if err != nil {
		return nil, err
	}
	return d.findOne(id)
}

// Find on file by its ID
func (d *DocDB) FindOne(id int) (*os.File, error) {
	doc, err := d.findOne(id)
	if err != nil {
		return nil, err
	}
	return d.docToFile(doc)
}

func (d *DocDB) findOne(id int) (*Doc, error) {
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
	var doc Doc
	if err = rows.Scan(&doc.ID, &doc.Name, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
		return nil, err
	}
	return &doc, nil
}

// Find many files. Useful for pagination.
// Use `DocCols` to select the right `ColName`
func (d *DocDB) FindMany(count int, offset int, orderCol ColName, ascending bool) ([]*os.File, error) {
	docs, err := d.findMany(count, offset, orderCol, ascending)
	if err != nil {
		return nil, err
	}
	return d.docsToFiles(docs)
}

func (d *DocDB) findMany(count int, offset int, orderCol ColName, ascending bool) ([]*Doc, error) {
	order := "DESC"
	if ascending {
		order = "ASC"
	}
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
// Use `DocCols` to select the right `ColName`
func (d *DocDB) FindAll(orderCol ColName, ascending bool) ([]*os.File, error) {
	docs, err := d.findAll(orderCol, ascending)
	if err != nil {
		return nil, err
	}
	return d.docsToFiles(docs)
}

func (d *DocDB) findAll(orderCol ColName, ascending bool) ([]*Doc, error) {
	order := "DESC"
	if ascending {
		order = "ASC"
	}
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
