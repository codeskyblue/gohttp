/* preview */

var React = require('react');
var ReactDOM = require('react-dom');

var Preview = require('./Preview.jsx');


ReactDOM.render(
	<div>
		<h3>Simple HTTP File Server</h3>
		<Preview/>
	</div>,
	document.getElementById('content')
)
