var React = require('react')
var {ButtonToolbar, DropdownButton, Button, 
	MenuItem, Modal} = require('react-bootstrap');
var humanize = require('humanize')
var urljoin = require('url-join');

var Icon = require('./Icon.jsx');
var FileIcon = require('./FileIcon.jsx');


var FileItem = React.createClass({
	getInitialState: function(){
		return {show: false};
	},
	render: function(){
		var fileType = this.props.data.type;
		var fileIcon;
		if (fileType == "directosry"){
			fileIcon = <Icon name="folder-open" style={{color: "#3366cc"}}/>
		} else {
			fileIcon = <FileIcon fileType={this.props.data.type} fileName={this.props.data.name}/>
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
		} else if (this.props.data.type == 'directory') {
			ctrlButtons.push(
				<Button key="zip-download"
					bsSize="xsmall" href={urljoin('/$zip', link)}>
					Download Zip <i className="fa fa-download"/>
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
			            				<img alt='qrcode' src={'/$qrcode?text='+encodeURI(urljoin(location.href, this.props.data.name))} />
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