var React = require('react');
var ReactDOM = require('react-dom');

var GithubButton = require('./github.jsx');
var Content = require('./content.jsx')

// ReactDOM.render(
// 	<GithubButton/>,
// 	document.getElementById("nav-right-bar")
// )

ReactDOM.render(
	<Content/>,
	document.getElementById('content')
)