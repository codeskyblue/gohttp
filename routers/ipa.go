package routers

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"

	goplist "github.com/DHowett/go-plist"
	"gopkg.in/macaron.v1"
)

func NewPlistHandler(rootDir string) macaron.Handler {
	return func(r *http.Request, w http.ResponseWriter, ctx *macaron.Context) {
		relpath := ctx.Params("*")
		if filepath.Ext(relpath) == ".plist" {
			relpath = relpath[0:len(relpath)-6] + ".ipa"
		}
		abspath := filepath.Join(rootDir, relpath)

		plinfo, err := parseIPA(abspath)
		if err != nil {
			log.Println(err)
			ctx.Error(500, err.Error())
			return
		}
		filepath.Ext(relpath)
		ipaURL := url.URL{
			Scheme: "https",
			Host:   r.Host,
			Path:   relpath,
		}
		imgURL := url.URL{
			Scheme: "https",
			Host:   r.Host,
			Path:   filepath.Join("/$ipaicon", relpath),
		}
		data, err := generateDownloadPlist(ipaURL.String(), imgURL.String(), plinfo)
		if err != nil {
			ctx.Error(500, err.Error())
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		w.Write(data)
	}
}

func NewIpaIconHandler(rootDir string) macaron.Handler {
	return func(w http.ResponseWriter, ctx *macaron.Context) {
		relpath := ctx.Params("*")
		abspath := filepath.Join(rootDir, relpath)
		data, err := parseIpaIcon(abspath)
		if err != nil {
			ctx.Error(404, err.Error())
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(data)
	}
}

func parseIpaIcon(path string) (data []byte, err error) {
	iconre := regexp.MustCompile(`(?i)^Payload/[^/]*/icon\.png$`)
	r, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer r.Close()

	var zfile *zip.File
	for _, file := range r.File {
		if iconre.MatchString(file.Name) {
			zfile = file
			break
		}
	}
	if zfile == nil {
		err = errors.New("icon.png file not found")
		return
	}
	plreader, err := zfile.Open()
	if err != nil {
		return
	}
	defer plreader.Close()
	return ioutil.ReadAll(plreader)
}

func IPAHandler(ctx *macaron.Context) {
	ctx.HTML(200, "ipa", nil)
}

type plistBundle struct {
	CFBundleIdentifier  string `plist:"CFBundleIdentifier"`
	CFBundleVersion     string `plist:"CFBundleVersion"`
	CFBundleDisplayName string `plist:"CFBundleDisplayName"`
}

func parseIPA(path string) (plinfo *plistBundle, err error) {
	plistre := regexp.MustCompile(`^Payload/[^/]*/Info\.plist$`)
	r, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer r.Close()

	var plfile *zip.File
	for _, file := range r.File {
		if plistre.MatchString(file.Name) {
			plfile = file
			break
		}
	}
	if plfile == nil {
		err = errors.New("Info.plist file not found")
		return
	}
	plreader, err := plfile.Open()
	if err != nil {
		return
	}
	defer plreader.Close()
	buf := make([]byte, plfile.FileInfo().Size())
	_, err = io.ReadFull(plreader, buf)
	if err != nil {
		return
	}
	dec := goplist.NewDecoder(bytes.NewReader(buf))
	plinfo = new(plistBundle)
	err = dec.Decode(plinfo)
	return
}

// ref: https://gist.github.com/frischmilch/b15d81eabb67925642bd#file_manifest.plist
type plAsset struct {
	Kind string `plist:"kind"`
	URL  string `plist:"url"`
}
type plItem struct {
	Assets   []*plAsset `plist:"assets"`
	Metadata struct {
		BundleIdentifier string `plist:"bundle-identifier"`
		BundleVersion    string `plist:"bundle-version"`
		Kind             string `plist:"kind"`
		Title            string `plist:"title"`
	} `plist:"metadata"`
}
type downloadPlist struct {
	Items []*plItem `plist:"items"`
}

func generateDownloadPlist(ipaUrl, imgUrl string, plinfo *plistBundle) ([]byte, error) {
	dp := new(downloadPlist)
	item := new(plItem)
	item.Assets = append(item.Assets, &plAsset{
		Kind: "software-package",
		URL:  ipaUrl,
	}, &plAsset{
		Kind: "display-image",
		URL:  imgUrl,
	})

	item.Metadata.Kind = "software"

	item.Metadata.BundleIdentifier = plinfo.CFBundleIdentifier
	item.Metadata.BundleVersion = plinfo.CFBundleVersion
	item.Metadata.Title = plinfo.CFBundleDisplayName

	dp.Items = append(dp.Items, item)
	data, err := goplist.MarshalIndent(dp, goplist.XMLFormat, "    ")
	// fmt.Println(string(data))
	// fmt.Println(err)
	return data, err
}
