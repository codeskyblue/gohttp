var React = require('react');
var PathBreadcrumb = require('./PathBreadcrumb.jsx');
var FilePreview = require('./FilePreview.jsx');
var path = require('path');


var Preview = React.createClass({
	getInitialState: function() {
		return {
			pathname: location.pathname,
			content: "loading ..."
		};
	},
	componentDidMount: function() {
		var ok = FilePreview.canPreview(this.state.pathname);
		ok && $.ajax(this.state.pathname, {
			method: 'GET',
			dataType: 'text',
			success: function(res){
				this.setState({content: res})
			}.bind(this),
			error: function(res){
				console.log(res)
				this.setState({content: res})
			}.bind(this)
		})
		if (!ok){
			this.setState({content: "File type not supported preview"})
		}
	},

	render: function() {
		var pathname = location.pathname;
		return (
			<div>
				<PathBreadcrumb data={pathname} />
				<FilePreview fileName={path.basename(pathname)} content={this.state.content}/>
			</div>
		);
	}

});

module.exports = Preview;