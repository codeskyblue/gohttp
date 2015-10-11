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
  },
  {
    name: ".gitignore",
    type: "file",
    size: 13322,
    mtime: 122123123844
  }
]

var FileList = React.createClass({
  render: function(){
    var that = this;
    var filterData = this.props.data.filter(function(item){
      var isHidden = (item.name.substr(0, 1) == '.')
      if (isHidden && !that.props.showHidden){
        return false;
      }
      return true;
    })
    var fileItems = filterData.map(function(item){
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
    var itemPaths = [];
    var currPath = '/';
    paths.forEach(function(name){
      var newPath = path.join(currPath, name);
      currPath = newPath;
      itemPaths.push({
        directory: currPath,
        name: name
      })
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
    return {
      data: fileList,
      hidden: false,
    }
  },
  render: function(){
    return (
      <Row>
        <Col md={12}>
          <PathBreadcrumb/>
        </Col>
        <Col md={12}>
          <Table striped bordered condensed hover>
            <thead>
              <tr>
                <td colSpan={5}>
                  <ButtonToolbar>
                    <Button bsSize="xsmall">
                      Upload <i className="fa fa-upload"/>
                    </Button>
                    <Button bsSize="xsmall"　onClick={
                      ()=>this.setState({hidden: !this.state.hidden})
                    }>
                      Show Hidden　<input type="checkbox" checked={this.state.hidden}/>
                    </Button>
                  </ButtonToolbar>
                </td>
              </tr>
              <tr>
                <th>
                  <Button bsSize="xsmall" bsStyle="link" href="..">
                    <i className="fa fa-level-up"/>
                  </Button>
                </th>
                <th className="table-name">Name</th>
                <th>Size</th>
                <th>Control</th>
                <th>Modity Time</th>
              </tr>
            </thead>
              
            <FileList data={this.state.data} showHidden={this.state.hidden}/>
          </Table>
        </Col>
      </Row>
    )
  }
})

module.exports = Explorer