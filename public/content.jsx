'use strict'

var React = require('react');
var Explorer = require('./explorer.jsx')

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