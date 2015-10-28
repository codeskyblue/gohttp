package routers

import "testing"

func TestGenerateDownloadPlist(t *testing.T) {
	p := new(plistBundle)
	p.CFBundleDisplayName = "tars"
	p.CFBundleIdentifier = "hahaha"
	p.CFBundleVersion = "0.0.1"

	generateDownloadPlist("https://^_^", "http://this is image", p)
}
