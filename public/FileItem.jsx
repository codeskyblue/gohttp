var React = require('react')
var {ButtonToolbar, DropdownButton, Button, 
	MenuItem, Modal} = require('react-bootstrap');
var humanize = require('humanize')
var urljoin = require('url-join');
var path = require('path')

var Icon = require('./Icon.jsx');
var FileIcon = require('./FileIcon.jsx');
var FilePreview = require('./FilePreview.jsx');

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

		var qrcodeLink = urljoin(location.href, this.props.data.name);
		if (path.extname(this.props.data.name) == '.ipa'){
			qrcodeLink = urljoin('https://'+location.host, '$ipa', location.pathname, this.props.data.name)
		}

		var link = urljoin(location.pathname, this.props.data.name);
		var fileLink = FilePreview.canPreview(this.props.data.name) ? link+'?preview=true' : link;
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
				<td className="hidden-xs text-center">
					{fileIcon}
				</td>
				<td>
					<a onClick={(e)=>this.props.onDirectoryChange && this.props.onDirectoryChange(link, e)} 
						href={fileLink}>{this.props.data.name}</a>
				</td>
				<td className="hidden-xs">
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
			            		<Modal.Title className="text-center">
			            			<a href={qrcodeLink}>{this.props.data.name}</a>
			            		</Modal.Title>
			            		<Modal.Body>
			            			<div className="text-center">
			            				<img alt='qrcode' src={'/$qrcode?text='+encodeURI(qrcodeLink)} />
			            			</div>
			            		</Modal.Body>
			            		<Modal.Footer>
			            			<Button onClick={close}>Close</Button>
			          			</Modal.Footer>
			          		</Modal.Header>
						</Modal>
					</div>
				</td>
				<td className="hidden-xs">{humanize.filesize(this.props.data.size)}</td>
				<td className="hidden-xs">
					{humanize.date('Y-m-d H:i:s', this.props.data.mtime)}
				</td>
			</tr>
		)
	}
})

module.exports = FileItem;
