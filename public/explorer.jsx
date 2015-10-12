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
var _ = require('underscore');
var path = require('path');
var urljoin = require('url-join');

var FileItem = require('./FileItem.jsx')
var PathBreadcrumb = require('./PathBreadcrumb.jsx')

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
      if (item.type == 'directory'){
        return (
          <FileItem key={item.name} data={item} onDirectoryChange={that.props.onDirectoryChange} />
        )
      } else {
        return (
          <FileItem key={item.name} data={item} />
        )
      }
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
    return {
      data: [],
      hidden: false,
      pathname: location.pathname,
    }
  },
  loadFilesFromServer: function(){
    // console.log(this.state)
    $.ajax({
      url: this.state.pathname+"?format=json",
      dataType: 'json',
      success: function(data){
        data = _.sortBy(data, function(item){
          return item.type+':'+item.name;
        })
        // console.log(data)
        this.setState({data: data})
      }.bind(this),
      error: function(xhr, status, err){
        console.log(status, err)
      }
    })
  },
  componentDidMount: function() {
    this.loadFilesFromServer();
  },
  changePath: function(newPath, e){
    e.preventDefault()
    // this.setState({pathname: newPath}) // WHY(ssx): setState not set immedietly
    this.state.pathname = newPath;
    window.history.pushState("", "", newPath);
    this.loadFilesFromServer();
  },
  render: function(){
    var key = 'hello'
    return (
      <Row>
        <Col md={12}>
          <PathBreadcrumb data={this.state.pathname} onClick={this.changePath} />
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
                  <Button bsSize="xsmall" href={urljoin(location.pathname, '..')}>
                    <i className="fa fa-level-up"/>
                  </Button>
                </th>
                <th className="table-name">Name</th>
                <th>Size</th>
                <th>Control</th>
                <th>Modity Time</th>
              </tr>
            </thead>
              
            <FileList data={this.state.data} showHidden={this.state.hidden} onDirectoryChange={this.changePath} />
            <tfoot>
              <tr>
                <td colSpan={5}>README.md # readme not finished
                </td>
              </tr>
            </tfoot>
          </Table>
        </Col>
      </Row>
    )
  }
})

module.exports = Explorer
