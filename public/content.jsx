'use strict'

var React = require('react');
var FileTable = require('./file-table.jsx')

var Content = React.createClass({
	render: function(){
		return (
			<div>
				<h3>{location.pathname}</h3>
				<FileTable/>
			</div>
		)
	}
})

module.exports = Content;