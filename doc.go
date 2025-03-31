package docdb

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Doc struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	filesPath string
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
func (d *Doc) Open() (*os.File, error) {
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
