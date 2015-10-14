var React = require('react');
var markdown = require('markdown').markdown;

var Markdown = React.createClass({

	render: function() {
		return (
			<div
				style={this.props.style} 
				dangerouslySetInnerHTML={{
            	__html: markdown.toHTML(this.props.text)
          	}}/>
		);
	}
});

module.exports = Markdown;