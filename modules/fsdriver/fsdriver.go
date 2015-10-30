package fsdriver

import (
	"errors"
	"io"
	"os"

	goftpserv "github.com/goftp/server"
)

var (
	ErrNotDirectory        = errors.New("Not a directory")
	ErrNotFile             = errors.New("Not a file")
	ErrDirHasSameName      = errors.New("A dir has the same name")
	ErrAppendFileNotExists = errors.New("Append data but file not exsit")
)

type FSDriver interface {
	Stat(path string) (os.FileInfo, error)
	ListDir(dir string) ([]os.FileInfo, error)
	Rename(oldpath string, newpath string) error
	Mkdir(dir string) error
	DeleteDir(dir string) error
	DeleteFile(file string) error

	// args: path, offset
	GetFile(path string, offset int64) (int64, io.ReadCloser, error)
	// args: path, reader, isAppend
	PutFile(path string, rd io.Reader, isAppend bool) (int64, error)
}

type GoftpDriverAdapter struct {
	goftpserv.Driver
}

func (d *GoftpDriverAdapter) ListDir(path string) ([]os.FileInfo, error) {
	sfis, err := d.DirContents(path)
	if err != nil {
		return nil, err
	}
	fis := make([]os.FileInfo, 0, len(sfis))
	for _, fi := range sfis {
		fis = append(fis, fi)
	}
	return fis, nil
	// return []os.FileInfo(fis), err
}

func (d *GoftpDriverAdapter) Stat(path string) (os.FileInfo, error) {
	ffi, err := d.Driver.Stat(path)
	return os.FileInfo(ffi), err
}
func (d *GoftpDriverAdapter) Mkdir(path string) error {
	return d.MakeDir(path)
}
