var humanize = require('humanize')
var React = require('react')
var Button = require('react-bootstrap').Button;
var ButtonToolbar = require('react-bootstrap').ButtonToolbar;
var DropdownButton = require('react-bootstrap').DropdownButton;
var MenuItem = require('react-bootstrap').MenuItem;
var Modal = require('react-bootstrap').Modal;
var urljoin = require('url-join');

var FileItem = React.createClass({
	getInitialState: function(){
		return {show: false};
	},
	render: function(){
		var fileType = this.props.data.type;
		var fileIcon;
		if (fileType == "directory"){
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
		var link = urljoin(location.pathname, this.props.data.name);

		var ctrlButtons = [];
		if (this.props.data.type == 'file'){
			ctrlButtons.push(
				<Button key="download"
					bsSize="xsmall" href={link+'?download=true'}>
					Download <i className="fa fa-download"/>
				</Button>,
				<Button key="qrcode" bsSize="xsmall" onClick={open}>
					QRCode <i className="fa fa-qrcode"/>
				</Button>
			)
		}
		return (
			<tr>
				<td className="text-center">
					{fileIcon}
				</td>
				<td>
					<a onClick={(e)=>this.props.onDirectoryChange(link, e)} href={link}>{this.props.data.name}</a>
				</td>
				<td>{humanize.filesize(this.props.data.size)}</td>
				<td>
					<div>
						<ButtonToolbar>
							{ctrlButtons}
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
			            				<img alt='qrcode' src={'/_qr?text='+urljoin(location.href, this.props.data.name)} />
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