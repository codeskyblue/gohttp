# HTTP Watcher

Focus on design and code. Hit F5 when needed for you, automatically.

No copy and paste javascript, http-watcher add the reload hook automatically.

>
HTTP file Server + HTTP proxy + Directory Watcher: Wacher directory for change, automatically reload connected bowsers. Works for both static and dynamic web project.

### build

```sh
  go build
```

### Usage

```sh
http-watcher args  # accepted args list below
```
```sh
  -command="": Command to run before reload browser, useful for preprocess, like compile scss. The files been chaneged, along with event type are pass as arguments
  -ignores="": Ignored file pattens, seprated by ',', used to ignore the filesystem events of some files
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

### The CPU usage is high

http-watcher is currently polling filesystem for event. When the directory is large, it may eat CPU
add more ignore pattens to filter files will make it lower

### TODO
Use kqueue for OS X, inotify for Windows for better performance
