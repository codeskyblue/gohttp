var React = require('react');

var Icon = React.createClass({

	render: function() {
		var className = ["fa"];
		if (this.props.name){
			className.push("fa-"+this.props.name)
		}
		if (this.props.size){
			className.push("fa-"+this.props.size)
		}
		return (
			<i className={className.join(' ')} style={this.props.style}/>
		);
	}

});

module.exports = Icon;