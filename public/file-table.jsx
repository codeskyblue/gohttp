'use strict'

var React = require('react')
var Table = require('react-bootstrap').Table;
var Button = require('react-bootstrap').Button;
var ButtonToolbar = require('react-bootstrap').ButtonToolbar;
var DropdownButton = require('react-bootstrap').DropdownButton;
var MenuItem = require('react-bootstrap').MenuItem;

var FileTable = React.createClass({
	render: function(){
		return (
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
				<tbody>
					<tr>
						<td>
							<i className="fa fa-toggle-left"/>
						</td>
						<td colSpan="4">
							<a href="../">Up directory</a>
						</td>
					</tr>
					<tr>
						<td><i className="fa fa-folder-open"/></td>
						<td>Jacob</td>
						<td>10M</td>
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
							3 days ago
						</td>
					</tr>
				</tbody>
			</Table>
		)
	}
})

module.exports = FileTable