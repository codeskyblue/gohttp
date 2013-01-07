# HTTP Watcher

HTTP file Server + File Watcher = HTTP Watcher => hit F5 for you

>
Wacher file for change, automatically reload connected bowser.

### build

```sh
  go build file-watcher.go
```

### Usage

```sh
Usage of ./http-watcher:
  -command="": Command to run before reload browser, useful for preprocess, like compile scss
  -ignores="": Ignored file pattens, seprated by ','
  -port=8000: Which port to listen
  -private=false: Only listen on lookback interface
  -root=".": Directory root been watched
```
