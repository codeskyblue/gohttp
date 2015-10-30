package fsdriver

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type PosixFSDriver struct {
	RootDir string
}

func (p *PosixFSDriver) IsDir(path string) bool {
	fi, err := p.Stat(path)
	return err == nil && fi.IsDir()
}

func (p *PosixFSDriver) ftp2fs(path string) string {
	return filepath.Join(p.RootDir, path)
}

func (p *PosixFSDriver) Stat(path string) (os.FileInfo, error) {
	return os.Lstat(p.ftp2fs(path))
}

func (p *PosixFSDriver) ListDir(path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(p.ftp2fs(path))
}

func (p *PosixFSDriver) Rename(oldpath, newpath string) error {
	return os.Rename(p.ftp2fs(oldpath), p.ftp2fs(newpath))
}

func (p *PosixFSDriver) Mkdir(path string) error {
	return os.Mkdir(p.ftp2fs(path), 0755)
}

func (p *PosixFSDriver) DeleteDir(path string) error {
	fi, err := p.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return ErrNotDirectory
	}
	return os.Remove(p.ftp2fs(path))
}

func (p *PosixFSDriver) DeleteFile(path string) error {
	fi, err := p.Stat(path)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return ErrNotFile
	}
	return os.Remove(p.ftp2fs(path))
}

func (p *PosixFSDriver) GetFile(path string, offset int64) (size int64, rc io.ReadCloser, err error) {
	fd, err := os.Open(p.ftp2fs(path))
	if err != nil {
		return
	}
	finfo, err := fd.Stat()
	if err != nil {
		fd.Close()
		return
	}
	_, err = fd.Seek(offset, os.SEEK_SET)
	if err != nil {
		return
	}
	return finfo.Size(), fd, nil
}

func (p *PosixFSDriver) PutFile(path string, data io.Reader, appendData bool) (int64, error) {
	fspath := p.ftp2fs(path)

	var exists bool
	f, err := p.Stat(path)
	if err == nil {
		exists = true
		if f.IsDir() {
			return 0, ErrDirHasSameName
		}
	} else {
		if os.IsNotExist(err) {
			exists = false
		} else {
			return 0, fmt.Errorf("Put File error: %v", err)
		}
	}

	if !appendData {
		if exists {
			err = p.DeleteFile(path)
			if err != nil {
				return 0, err
			}
		}
		f, err := os.Create(fspath)
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

	if !exists {
		return 0, ErrAppendFileNotExists
	}

	file, err := os.OpenFile(fspath, os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	_, err = file.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}

	bytes, err := io.Copy(file, data)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
