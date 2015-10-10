'use strict'

var React = require('react');
var ReactDOM = require('react-dom');
var rbs = require('react-bootstrap');

var NavItem = rbs.NavItem,
	Nav = rbs.Nav,
	Button = rbs.Button,
	Modal = rbs.Modal,
	ButtonToolbar = rbs.ButtonToolbar;

var githubLink = "https://github.com/codeskyblue/file-server";
var GithubButton = React.createClass({
	render: function(){
		return (
			<NavItem href={githubLink}>
				<i className="fa fa-github"/> GITHUB
			</NavItem>
		)
	}
})

var QRCode = React.createClass({
	getInitialState(){
		return {showModal: false};
	},
	close(){
		this.setState({showModal: false});
	},
	open(){
		this.setState({showModal: true});
	},
	render: function(){
		return (
			<span>
				<Button onClick={this.open} bsSize="xsmall">
					二维码 <i className="fa fa-qrcode"/>
				</Button>
				<Modal bsSize="small" show={this.state.showModal} onHide={this.close}>
					<Modal.Header closeButton>
						<Modal.Title>.gitignore</Modal.Title>
					</Modal.Header>
					<Modal.Body>
						<img src="/_qr?text='hallo'"/>
						<p>打开手机扫一扫，对准二维码扫描即可下载地址。</p>
					</Modal.Body>
				</Modal>
			</span>
		)
	}
})

var AccessLink = React.createClass({
	render: function(){
		return (
			<ButtonToolbar>
				<QRCode/>
				<QRCode/>
			</ButtonToolbar>
		)
	}
})

module.exports = GithubButton

// ReactDOM.render(
// 	<GithubButton/>,
// 	document.getElementById("nav-right-bar")
// )

// var qrcodes = document.getElementsByClassName("qrcode")

// Array.prototype.forEach.call(qrcodes, function(mountNode){
// 	ReactDOM.render(
// 		<AccessLink/>, mountNode
// 	)
// })
