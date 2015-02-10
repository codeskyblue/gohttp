package main

const (
	MODIFY   = "MODIFY"
	ADD      = "ADD"
	REMOVE   = "REMOVE"
	DIR_HTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Directory Listing: {{.dir}}</title>
    <style>
      body { font: 14px/1.4 Monospace; }
      #page-wrap { margin: 0 auto; width: 800px; }
     	@media (max-width: 1000px) {
        #page-wrap { width: auto; padding: 0 4px; }
        .mtime { font-size: 9px; white-space: nowrap;}
      	table, tr { width: 100%;}
      	.name {
	         max-width: 450px;
           display: inline-block;
	         text-overflow: ellipsis;
           overflow: hidden;
	      }
      }
      table { width: 100%; }
      caption {
          font-weight: bold;
          font-size: 18px;
          margin: 20px;
      }
      thead {
          font-size: 15px;
          background: #DFF0D8;
      }
      th, td { padding: 5px 2px; }
      tr:nth-child(2n) { background: #eee;n }
      tr:nth-child(2n) td { background: #eee; }
      #footer {
          margin: 20px 0;
          text-align: right;
          font-size: 11px;
          color: #888;
      }
      #footer a { color: #555; }
    </style>
  </head>
  <body>
    <div id="page-wrap">
      <table cellspacing=0>
        <caption>Directory List:  {{.dir}}</caption>
        <thead>
          <th>File</th>
          <th>Size</th>
          <th>Last Modified</th>
        </thead>
        {{range .files}}
        <tr>
          <td class="name"><a href="{{.href}}">{{.name}}</a></td>
          <td>{{ .size }}</td>
          <td class="mtime">{{ .mtime }}</td>
        </tr>
        {{end}}
      </table>
			<div id="footer">
        <a href="https://github.com/shenfeng/http-watcher">http-watcher</a>,
        write by <a href="http://shenfeng.me">Feng Shen</a> in golang,<a href="/_d/doc">doc</a>
      </div>
    </div>
  </body>
</html>`
	HELP_HTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>HTTP watcher Documentation</title>
    <style>
      body { width: 800px; margin: 0 auto; }
	    @media (max-width: 1000px) {
        body { width: auto; padding: 0 4px; }
      }
      .mesg { background: #fff1a8; padding: 6px 2px; font: bold 15px monospace; }
      .note {
      	  background: #ffffcc;
      	  font-family: monospace;
      	  padding: 4px;
      }
      pre { white-space: pre-wrap;}
      h1 { text-align: center;}
      ul { padding: 0; list-style: none; }
      li { padding: 4px; margin: 4px 0; }
      #footer {
          margin: 20px 0;
          text-align: right;
          font-size: 11px;
          color: #888;
      }
      #footer .doc { float: left; }
      #footer a { color: #555; }
      .ignores li { padding: 0; margin: 4px; }
      pre {
         font-family: monospace;
         font-size: 15px;
         line-height: 1.4;
      }
    </style>
  </head>
  <body>
    <h1>HTTP Watcher Documentation</h1>
    {{if .error}}
      <p class="mesg">ERROR: {{.error}}</p>
    {{end}}

    <h3>Directory been watched for changed</h3>
    <p class="note">{{.dir}}</p>
    <div>
      <p class="note">Ignore file pattens: </p>
      <ol class="ignores">
		{{range .ignores}}<li>{{.}}</li>{{end}}
      </ol>
    </div>
    <h3>Visit (automatically reload when file changes detected)</h3>
    <ul>
      {{range .hosts}}
        <li class="note"><a href="http://{{.}}/">http://{{.}}/</a></li>
      {{end}}
    </ul>
    <h3>Command Help</h3>
    <pre>http-watcher -h</pre>
    </div>
    <div id="footer">
       <a href="https://github.com/shenfeng/http-watcher">http-watcher</a>,
       write by <a href="http://shenfeng.me">Feng Shen</a> in golang,<a href="/_d/doc">doc</a>
    </div>
  </body>
</html>`
)
