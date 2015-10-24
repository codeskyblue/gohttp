var React = require('react');

var PreviewImage = React.createClass({

	render: function() {
		return (
			<div style={this.props.style}>
				<img src={this.props.fileName}/>
			</div>
		);
	}
});

module.exports = PreviewImage;
