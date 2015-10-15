var React = require('react');
var Markdown = require('./Markdown.jsx');
var request = require('superagent')
var path = require('path');


// FIXME(ssx): Here got a duplicated request problem.
var FilePreview = React.createClass({
	getInitialState: function() {
		return {
			fileName: "",
			content: "",
		};
	},
	updatePreviewFile: function(fileName){
		console.log("preview:", fileName)
		var that = this;
		request.get(fileName)
			.end(function(err, res){
				console.log("Request:", fileName)
				if (err){
					console.log(err);
					return
				}
				that.setState({
					fileName: fileName,
					content: res.text,
				})
			})
	},
	componentDidMount: function() {
		this.updatePreviewFile(this.props.fileName)
	},
	shouldComponentUpdate: function(nextProps, nextState) {
		if(! nextProps){
			return true;
		}
		return this.state.fileName !== nextProps.fileName
	},
	componentWillUpdate: function(nextProps, nextState) {
		this.updatePreviewFile(this.props.fileName)
	},
	render: function() {
		var fileName = this.state.fileName || "";
		console.log("Rendered", fileName)
		var ext = path.extname(fileName.toLowerCase());
		switch(ext){
		case "":
		case ".txt":
			return <pre>{this.state.content}</pre>;
		case ".md":
			return <Markdown text={this.state.content} style={{margin: '0px 15px'}}/>;
		default:
			return <span><b>Not supported file type</b> <i>{fileName}</i></span>;
		}
	}
});

module.exports = FilePreview;