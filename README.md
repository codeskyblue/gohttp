# HTTP Watcher

HTTP file Server + File Watcher = HTTP Watcher => hit F5 for you

>
Wacher file for change, automatically reload connected bowser. Just start the server by run `http-watcher`, happy design your site, http-watcher did the rest for you
One source file, one binary, works on everywhare golang works

### build

```sh
  go build file-watcher.go
```

### Command line args

```sh
Usage of ./http-watcher:
  -command="": Command to run before reload browser, useful for preprocess, like compile scss. The files been chaneged, along with event type are pass as arguments
  -ignores="": Ignored file pattens, seprated by ','
  -port=8000: Which port to listen
  -private=false: Only listen on lookback interface
  -root=".": Directory root been watched
```


### Static Web Site

Just run `http-watcher` on the HTML root directory, then happy coding.

### For dynamic web site, like Clojure, go web project
Copy and paste to your HTML: (make sure http-watcher is running :) )
```html
<script src="http://127.0.0.1:8000/_d/js"></script>
```
