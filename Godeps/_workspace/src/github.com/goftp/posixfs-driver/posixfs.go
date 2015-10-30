package posixfs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/goftp/server"
)

type BasicFileInfo struct {
	os.FileInfo
}

func (f *BasicFileInfo) Owner() string {
	return "root"
}

func (f *BasicFileInfo) Group() string {
	return "root"
}

type PosixFSDriver struct {
	RootPath string
}

func (driver *PosixFSDriver) ChangeDir(path string) error {
	rPath := filepath.Join(driver.RootPath, path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	return errors.New("Not a dir")
}

func (driver *PosixFSDriver) Stat(path string) (server.FileInfo, error) {
	basepath := filepath.Join(driver.RootPath, path)
	rPath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	f, err := os.Lstat(rPath)
	if err != nil {
		return nil, err
	}
	return &BasicFileInfo{f}, nil
}

func (driver *PosixFSDriver) DirContents(path string) ([]server.FileInfo, error) {
	basepath := filepath.Join(driver.RootPath, path)
	fis, err := ioutil.ReadDir(basepath)
	if err != nil {
		return nil, err
	}

	files := make([]server.FileInfo, 0)
	for _, finfo := range fis {
		files = append(files, &BasicFileInfo{finfo})
	}
	return files, nil
}

func (driver *PosixFSDriver) DeleteDir(path string) error {
	rPath := filepath.Join(driver.RootPath, path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a directory")
}

func (driver *PosixFSDriver) DeleteFile(path string) error {
	rPath := filepath.Join(driver.RootPath, path)
	f, err := os.Lstat(rPath)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		return os.Remove(rPath)
	}
	return errors.New("Not a file")
}

func (driver *PosixFSDriver) Rename(fromPath string, toPath string) error {
	oldPath := filepath.Join(driver.RootPath, fromPath)
	newPath := filepath.Join(driver.RootPath, toPath)
	return os.Rename(oldPath, newPath)
}

func (driver *PosixFSDriver) MakeDir(path string) error {
	rPath := filepath.Join(driver.RootPath, path)
	return os.Mkdir(rPath, os.ModePerm)
}

func (driver *PosixFSDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	rPath := filepath.Join(driver.RootPath, path)
	f, err := os.Open(rPath)
	if err != nil {
		return 0, nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return 0, nil, err
	}

	f.Seek(offset, os.SEEK_SET)
	return info.Size(), f, nil
}

func (driver *PosixFSDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	rPath := filepath.Join(driver.RootPath, destPath)
	var isExist bool
	f, err := os.Lstat(rPath)
	if err == nil {
		isExist = true
		if f.IsDir() {
			return 0, errors.New("A dir has the same name")
		}
	} else {
		if os.IsNotExist(err) {
			isExist = false
		} else {
			return 0, errors.New(fmt.Sprintln("Put File error:", err))
		}
	}

	if !appendData {
		if isExist {
			err = os.Remove(rPath)
			if err != nil {
				return 0, err
			}
		}
		f, err := os.Create(rPath)
		if err != nil {
			return 0, err
		}
		defer f.Close()
		bytes, err := io.Copy(f, data)
		if err != nil {
			return 0, err
		}
		return bytes, nil
	}

	if !isExist {
		return 0, errors.New("Append data but file not exsit")
	}

	of, err := os.OpenFile(rPath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer of.Close()

	_, err = of.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(of, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}

func NewPosixFSFactory(rootDir string) *PosixFSFactory {
	return &PosixFSFactory{rootDir}
}

type PosixFSFactory struct {
	RootPath string
}

func (bdf *PosixFSFactory) NewDriver() (server.Driver, error) {
	return &PosixFSDriver{bdf.RootPath}, nil
}
