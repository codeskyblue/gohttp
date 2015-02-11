#
# Makefile
# hzsunshx, 2015-02-11 13:17
#

all:
	go-bindata templates public/...
	go build
	GOOS=windows GOARCH=386 go build


# vim:ft=make
#
