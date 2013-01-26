# HTTP Watcher

A server that automatically reload browsers when file changed, help developers focus on coding.

No copy and paste javascript code needed, just start `http-watcher`, that's all.

>
Web Server for Web developers! HTTP Watcher = HTTP file Server + HTTP proxy + Directory Watcher: automatically reload connected browsers when file changed, works for both static and dynamic web project.

### build

```sh
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

Start `http-watcher` with -root=$PROJECT_ROOT -port $PORT_NUMBER

where

* `PROJECT_ROOT` : the static web project's root directory, where `http-watcher` watch for filesystem events (MODIFY, ADD, REMOVE). Default: current directory
* `PORT_NUMBER` : which port `http-watcher` listens. Default: 8000

Now visit: [http://127.0.0.1:8000](http://127.0.0.1:800), if you take the default PORT_NUMBER

### Dynamic web site: Clojure, golang, Python, JAVA, etc project

Start `http-watcher` with -proxy=$PROXY_PORT -root=$PROJECT_ROOT -port $PORT_NUMBER

Where

* `PROXY_PORT` : the port the dynamic web project is listened on
* `PROJECT_ROOT` : the dynamic web project's root directory, where `http-watcher` watch for filesystem events (MODIFY, ADD, REMOVE). Default: current directory
* `PORT_NUMBER` : which port `http-watcher` listen. Default: 8000

Now visit: [http://127.0.0.1:8000](http://127.0.0.1:800), if you take the default PORT_NUMBER

also accept -ignores='regexp1,regexp2,regexp3' to ignore centain files

also accept -command=$SCRIPT_PATH to do some preprocessing before reload browsers. Take `preprocess` script as an example

`http-watcher` acts as a proxy in this configration

### HTTP file server, no filesystem monitoring

with arg : -monitor=false

### The CPU usage is high

http-watcher is currently polling filesystem for event. When the directory is large, it may eat CPU
Adding more ignore pattens to filter files will make it lower

