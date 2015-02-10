# File Static Server
This is file server writen by golang.

This is a short version of [http-watcher](https://github.com/shenfeng/http-watche)

http-watcher的简化版，除去了其他东西，只保留了文件服务器的东西。

### build

```sh
  go clone https://github.com/codeskyblue/file-server
  go build  # you may want to copy http-watcher binary to $PATH for easy use. prebuilt binary comming soon
```

### Usage

```sh
file-server ARGS  # acceptable args list below, -h to show them
```
```sh
  -port=8000: Which port to listen
  -private=false: Only listen on lookback interface, otherwise listen on all interface
  -root=".": the HTTP File Server's root directory
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

## Thanks
1. <https://github.com/shenfeng/http-watcher>