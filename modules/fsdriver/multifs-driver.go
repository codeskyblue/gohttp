package fsdriver

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/log"
	"gopkg.in/yaml.v2"
)

const (
	GOHTTPCFG   = ".gohttp.yml"
	MOUNT_QINIU = "qiniu"
)

type MountItem struct {
	Type   string
	Config map[string]string
}

type HtConfig struct {
	Mount MountItem `yaml:"mount"`
}

func readHtConfig(rd io.Reader) (*HtConfig, error) {
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	ht := &HtConfig{}
	err = yaml.Unmarshal(data, ht)
	return ht, err
}

func makeDriver(mcfg MountItem) FSDriver {
	if mcfg.Type == MOUNT_QINIU {
		dr, _ := NewQiniuDriver(
			mcfg.Config["access_key"],
			mcfg.Config["secret_key"],
			mcfg.Config["bucket"])
		return dr
	}
	return nil
}

type MultiFSDriver struct {
	pfsdrv *PosixFSDriver
}

func (p *MultiFSDriver) Stat(path string) (os.FileInfo, error) {
	rel, dr := p.getDriver(path)
	return dr.Stat(rel)
}

func (p *MultiFSDriver) ListDir(path string) (fis []os.FileInfo, err error) {
	rel, dr := p.getDriver(path)
	return dr.ListDir(rel)
}

// Rename folder multi driver is not supported
func (p *MultiFSDriver) Rename(oldpath string, newpath string) error {
	oldpath = filepath.Clean(oldpath)
	newpath = filepath.Clean(newpath)

	oldRel, oldDriver := p.getDriver(oldpath)
	oldBase := oldpath[0 : len(oldpath)-len(oldRel)]
	_, err := filepath.Rel(oldBase, newpath)

	if err == nil {
		return oldDriver.Rename(oldpath, newpath)
	} else {
		// oldDriver.Stat(oldRel)
		// FIXME(ssx): Check if oldRel isDir then return Error
		newRel, newDriver := p.getDriver(newpath)
		_, rc, err := oldDriver.GetFile(oldRel, 0)
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = newDriver.PutFile(newRel, rc, false)
		if err == nil {
			oldDriver.DeleteFile(oldRel)
		}
		return err
	}
}

func (p *MultiFSDriver) Mkdir(path string) error {
	rel, dr := p.getDriver(path)
	return dr.Mkdir(rel)
}

func (p *MultiFSDriver) DeleteDir(path string) error {
	rel, dr := p.getDriver(path)
	return dr.DeleteDir(rel)
}

func (p *MultiFSDriver) DeleteFile(path string) error {
	rel, dr := p.getDriver(path)
	return dr.DeleteFile(rel)
}

func (p *MultiFSDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	rel, dr := p.getDriver(path)
	return dr.GetFile(rel, offset)
}

func (p *MultiFSDriver) PutFile(path string, rd io.Reader, isAppend bool) (int64, error) {
	rel, dr := p.getDriver(path)
	return dr.PutFile(rel, rd, isAppend)
}

/*
	- Stat(path string) (os.FileInfo, error)
	- ListDir(dir string) ([]os.FileInfo, error)
	- Rename(oldpath string, newpath string) error
	- Mkdir(dir string) error
	- DeleteDir(dir string) error
	- DeleteFile(file string) error

	// args: path, offset
	- GetFile(path string, offset int64) (int64, io.ReadCloser, error)
	// args: path, reader, isAppend
	- PutFile(path string, rd io.Reader, isAppend bool) (int64, error)
*/

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	path = filepath.Clean(path)
	return strings.Split(path, "/")
}

func (mdr *MultiFSDriver) getDriver(path string) (rel string, driver FSDriver) {
	p := mdr.pfsdrv
	names := splitPath(path)
	stack := []string{}
	for idx, name := range names {
		stack = append(stack, name)
		tmppath := filepath.Join(stack...)
		if !p.IsDir(tmppath) {
			continue
		}
		cfgpath := filepath.Join(tmppath, GOHTTPCFG)
		_, err := p.Stat(cfgpath)
		if err != nil {
			continue
		}

		// check mount
		_, rd, err := p.GetFile(cfgpath, 0)
		if err != nil {
			continue
		}
		defer rd.Close()

		hcfg, err := readHtConfig(rd)
		if err != nil {
			log.Warn(err)
			continue
		}
		if hcfg.Mount.Type != "" {
			log.Debug(hcfg.Mount.Type, "/"+filepath.Join(names[idx+1:]...))
			return "/" + filepath.Join(names[idx+1:]...), makeDriver(hcfg.Mount)
		}
	}

	return path, p
}
