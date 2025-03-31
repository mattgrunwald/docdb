package docdb

import "database/sql"

type DocDB struct {
	db        *sql.DB
	filesPath string
}

func (d *DocDB) newDoc() Doc {
	return Doc{
		filesPath: d.filesPath,
	}
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

func (d *DocDB) Close() error {
	return d.db.Close()
}
