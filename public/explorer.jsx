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
var Modal = require('react-bootstrap').Modal;
var ReactBS = require('react-bootstrap');
var Row = ReactBS.Row,
  Col = ReactBS.Col,
  Breadcrumb = ReactBS.Breadcrumb,
  BreadcrumbItem = ReactBS.BreadcrumbItem;

var path = require('path')

var FileItem = require('./file-item.jsx')

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

var PathBreadcrumb = React.createClass({
  render: function(){
    var paths = location.pathname.split('/');
    paths.pop();
    paths.shift();
    console.log(paths);
    var itemPaths = [];
    var currPath = '/';
    paths.forEach(function(name){
      // var currPath = path.join(, 's')
      var newPath = path.join(currPath, name);
      currPath = newPath;
      itemPaths.push({
        directory: currPath,
        name: name
      })
      console.log(currPath)
    })
    var items = itemPaths.map(function(subPath){
      return (
        <BreadcrumbItem key={subPath.directory} href={subPath.directory}>
          {subPath.name}
        </BreadcrumbItem>
      )
    })
    return (
      <Breadcrumb>
        <BreadcrumbItem href="/">
          $
        </BreadcrumbItem>
        {items}
      </Breadcrumb>
    )
  }
})

var Explorer = React.createClass({
  getInitialState: function(){
    return {data: fileList}
  },
  render: function(){
    return (
      <Row>
        <Col md={3}>
          <ButtonToolbar>
            <Button bsSize="small">
              Up <i className="fa fa-level-up"/>
            </Button>
            <Button bsSize="small">
              Upload <i className="fa fa-upload"/>
            </Button>
          </ButtonToolbar>
        </Col>
        <Col md={9}>
          <PathBreadcrumb/>
        </Col>
        <Col md={12}>
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
        </Col>
      </Row>
    )
  }
})

module.exports = Explorer