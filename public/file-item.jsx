var humanize = require('humanize')
var React = require('react')
var Button = require('react-bootstrap').Button;
var ButtonToolbar = require('react-bootstrap').ButtonToolbar;
var DropdownButton = require('react-bootstrap').DropdownButton;
var MenuItem = require('react-bootstrap').MenuItem;
var Modal = require('react-bootstrap').Modal;


var FileItem = React.createClass({
	getInitialState: function(){
		return {show: false};
	},
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
		var that = this;
		var close = function(){
			that.setState({show: false})
		}
		var open = function(){
			that.setState({show: true})
		}
		return (
			<tr>
				<td>
					{fileIcon}
				</td>
				<td>
					<a href={this.props.data.name}>{this.props.data.name}</a>
				</td>
				<td>{humanize.filesize(this.props.data.size)}</td>
				<td>
					<div>
						<ButtonToolbar>
							<Button bsSize="xsmall" href={this.props.data.name}>
								Download <i className="fa fa-download"/>
							</Button>
							<Button bsSize="xsmall" onClick={open}>
								QRCode <i className="fa fa-qrcode"/>
							</Button>
						</ButtonToolbar>
						<Modal
							bsSize="small"
							show={this.state.show}
							onHide={close}
						>
							<Modal.Header closeButton>
            		<Modal.Title className="text-center">{this.props.data.name}</Modal.Title>
            		<Modal.Body>
            			<div className="text-center">
            				<img alt='qrcode' src={'/_qr?text='+location.href+this.props.data.name} />
            			</div>
            		</Modal.Body>
            		<Modal.Footer>
            			<Button onClick={close}>Close</Button>
          			</Modal.Footer>
          		</Modal.Header>
						</Modal>
					</div>
				</td>
				<td>
					{humanize.date('Y-m-d H:i:s', this.props.data.mtime)}
				</td>
			</tr>
		)
	}
})

module.exports = FileItem;