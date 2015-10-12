'use strict'

var React = require('react');
var Explorer = require('./Explorer.jsx')

var Content = React.createClass({
	render: function(){
		return (
			<div>
				<h3>Simple HTTP File Server</h3>
				<Explorer/>
			</div>
		)
	}
})

module.exports = Content;