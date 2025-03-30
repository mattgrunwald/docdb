package docdb

import (
	"database/sql"
	"os"
	"time"

	"github.com/mattgrunwald/docdb/col"
	"github.com/mattgrunwald/docdb/order"
)

type IDocDB interface {
	Insert(file *os.File) (*Doc, error)
	Update(id int, file *os.File) (*Doc, error)
	FindOne(id int) (*Doc, error)
	FindMany(count int, offset int, orderCol col.DocCol, order order.Order) ([]*Doc, error)
	FindAll(orderCol col.DocCol, order order.Order) ([]*Doc, error)
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
	filesPath string
}
