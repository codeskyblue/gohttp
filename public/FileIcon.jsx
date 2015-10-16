var React = require('react');
var path = require('path');
var Icon = require('./Icon.jsx');


var FindIcon = React.createClass({

	render: function() {
		if (this.props.name){
			return <Icon name={this.props.name}/>
		}

		var extname = path.extname(this.props.fileName || "").toLowerCase();
		if (this.props.fileType == "directory"){
			if (this.props.fileName == ".git"){
				return <Icon name="git-square"/>
			}
			return <Icon name="folder-open" style={{color: "#3366cc"}}/>
		}
		switch(extname) {
			case ".gif":
			case ".png":
			case ".jpg":
			case ".jpeg":
				return <Icon name="file-image-o"/>
			case ".zip":
				return <Icon name="file-zip-o"/>
			case ".apk":
				return <Icon name="android"/>
			case ".ipa":
				return <Icon name="apple"/>
			case ".exe":
				return <Icon name="windows"/>
			case ".txt":
				return <Icon name="file-text-o"/>
			default:
				return <Icon name="file-o"/>
				break;
		}
	}

});

module.exports = FindIcon;