package docdb

import (
	"database/sql"
	"os"
	"time"
)

type IDocDB interface {
	Insert(file *os.File) (*Doc, error)
	Update(id int, file *os.File) (*Doc, error)
	FindOne(id int) (*os.File, error)
	FindMany(count int, offset int, orderCol ColName, ascending bool) ([]*os.File, error)
	FindAll(orderCol ColName, ascending bool) ([]*os.File, error)
	Delete(id int) error
}

type DocDB struct {
	db        *sql.DB
	filesPath string
}

type Doc struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ColName string

type OrderCols struct {
	ID        ColName
	Name      ColName
	CreatedAt ColName
	UpdatedAt ColName
}

var DocCols = OrderCols{
	ID:        "id",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}
