'use strict'

// refs:
// https://github.com/taijinlee/humanize

var humanize = require('humanize')
var React = require('react')
var Table = require('react-bootstrap').Table;
var Button = require('react-bootstrap').Button;
var ButtonToolbar = require('react-bootstrap').ButtonToolbar;
var DropdownButton = require('react-bootstrap').DropdownButton;
var MenuItem = require('react-bootstrap').MenuItem;


var fileList = [
	{
		name: "hello.js",
		type: "directory",
		size: 1323003,
		mtime: 12283844
	},
	{
		name: "world.md",
		type: "file",
		size: 13322,
		mtime: 122123123844
	}
]

var FileItem = React.createClass({
	render: function(){
		var fileType = this.props.data.type;
		var fileIcon;
		console.log(fileType);
		if (fileType == "directory"){
			console.log("Is dir")
			fileIcon = <i className="fa fa-folder-open"/>;
		} else {
			fileIcon = <i className="fa fa-file-o"/>;
		}
		return (
			<tr>
				<td>
					{fileIcon}
				</td>
				<td>{this.props.data.name}</td>
				<td>{humanize.filesize(this.props.data.size)}</td>
				<td>
					<ButtonToolbar>
						<Button bsSize="xsmall">
							Download <i className="fa fa-download"/>
						</Button>
						<Button bsSize="xsmall">
							QRCode <i className="fa fa-qrcode"/>
						</Button>
					</ButtonToolbar>
				</td>
				<td>
					{humanize.date('Y-m-d H:i:s', this.props.data.mtime)}
				</td>
			</tr>
		)
	}
})

var FileList = React.createClass({
	render: function(){
		var fileItems = this.props.data.map(function(item){
			return (
				<FileItem key={item.name} data={item}/>
			)
		})
		return (
			<tbody>
				{fileItems}
			</tbody>
		)
	}
})

var Explorer = React.createClass({
	getInitialState: function(){
		return {data: fileList}
	},
	render: function(){
		return (
			<div>
				<ButtonToolbar>
					<Button bsSize="small">
						Up <i className="fa fa-level-up"/>
					</Button>
					<Button bsSize="small">
						Upload <i className="fa fa-upload"/>
					</Button>
					<Button bsStyle="link" bsSize="small">
						<span>{location.pathname}</span>
					</Button>
				</ButtonToolbar>
				<br/>
				<Table striped bordered condensed hover>
					<thead>
						<tr>
							<th>#</th>
							<th className="table-name">Name</th>
							<th>Size</th>
							<th>Control</th>
							<th>Modity Time</th>
						</tr>
					</thead>
						
					<FileList data={this.state.data} />
				</Table>
			</div>
		)
	}
})

module.exports = Explorer