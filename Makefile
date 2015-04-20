#
# Makefile
# hzsunshx, 2015-02-11 13:17
#

all:
	go-bindata data templates public/...
	GOOS=windows GOARCH=386 go build
	GOOS=linux GOARCH=386 go build -o fileserv-linux-386
	GOOS=linux GOARCH=amd64 go build -o fileserv-linux-amd64


# vim:ft=make
#
