var React = require('react');
var request = require('superagent')
var path = require('path');
var Markdown = require('./Markdown.jsx');


// FIXME(ssx): Here got a duplicated request problem.
var FilePreview = React.createClass({
	render: function() {
		var fileName = this.props.fileName || "";
		var ext = path.extname(fileName.toLowerCase());
		switch(ext){
		case "":
			return <span>{this.props.content}</span>;
		case ".txt":
			return <pre>{this.props.content}</pre>;
		case ".md":
			return <Markdown text={this.props.content} style={{margin: '0px 15px'}}/>;
		default:
			return <span><b>Not supported file type</b> <i>{fileName}</i></span>;
		}
	}
});

module.exports = FilePreview;