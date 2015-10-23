var React = require('react');
var request = require('superagent')
var path = require('path');
var Markdown = require('./Markdown.jsx');
var {Panel} = require('react-bootstrap');
var Highlight = require('react-highlight');


var highlightExtentions = [
	".txt",
	".js",
	".go",
	".json",
	".py",
	".html",
	".yml",
	".yaml",
	".c",
	".cpp",
	".h",
	".xml",
	".php",
	".sh",
	".jsx",
];

var FilePreview = React.createClass({
	render: function() {
		var fileName = this.props.fileName || "";
		var ext = path.extname(fileName.toLowerCase());
		var contentNode = null;
		switch(ext){
		case "":
			contentNode = <span>{this.props.content}</span>;
			break;
		case ".log":
			contentNode = <pre>{this.props.content}</pre>;
			break;
		case ".md":
			contentNode = <Markdown text={this.props.content} style={{margin: '0px 15px'}}/>;
			break;
		default:
			if (highlightExtentions.indexOf(ext) !== -1){
				var lang = ({
					js: 'javascript',
					yml: 'yaml'
				})[ext.substr(1)] || ext.substr(1);
				// FIXME(ssx): add line-numbers support, Don't know
				contentNode = <Highlight className={lang}>{this.props.content}</Highlight>
			} else {
				contentNode = <p><b>Prewview not supported file </b> <i>{fileName}</i></p>;
			}
			break;
		}
		
		return (
			<Panel className='panel-code-body' header={path.basename(fileName)}>
				{contentNode}
			</Panel>
		)
	}
});

FilePreview.canPreview = function(filename){
	var ext = path.extname(filename);
	return highlightExtentions.concat('.log', '.txt', '.md').indexOf(ext) !== -1;
}

module.exports = FilePreview;
