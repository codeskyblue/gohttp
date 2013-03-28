# HTTP Watcher

A server that automatically reload browsers when file changed, help developers focus on coding.

No copy and paste javascript code needed, just start `http-watcher`, that's all.

>
Web Server for Web developers! HTTP Watcher = HTTP file Server + HTTP proxy + Directory Watcher: automatically reload connected browsers when file changed, works for both static and dynamic web project.

### build

```sh
  # go get github.com/howeyc/fsnotify
  go build  # you may want to copy http-watcher binary to $PATH for easy use. prebuilt binary comming soon
```

### Usage

```sh
http-watcher args  # acceptable args list below, -h to show them
```
```sh
  -command="": Command to run before reload browser, useful for preprocess, like compile scss. The files been chaneged, along with event type are pass as arguments
  -ignores="": Ignored file pattens, seprated by ',', used to ignore the filesystem events of some files
  -monitor=true: Enable monitor filesystem event
  -port=8000: Which port to listen
  -private=false: Only listen on lookback interface, otherwise listen on all interface
  -proxy=0: Local dynamic site's port number, like 8080, HTTP watcher proxy it, automatically reload browsers when watched directory's file changed
  -root=".": Watched root directory for filesystem events, also the HTTP File Server's root directory
```

### HTML + JS + CSS (static web project)

```sh
http-watcher -port 8000 -root /your/code/root
```

### Dynamic web site: Clojure, golang, Python, JAVA

```sh
# your dynamic site listen on 9090
# http-watcher act as a proxy
http-watcher -port 8000 -root /your/code/root -proxy=9090 -ignores test/,classes
```
### HTTP file server, no filesystem monitoring

```sh
# like python -m SimpleHTTPServer, should handle concurrency better
http-watcher -monitor=false
```
