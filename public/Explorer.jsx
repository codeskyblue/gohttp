'use strict'

// refs:
// https://github.com/taijinlee/humanize

var React = require('react')
var {Row, Col, ButtonToolbar, MenuItem, Modal,
  DropdownButton, Breadcrumb, BreadcrumbItem, Button, Table} = require('react-bootstrap');

var _ = require('underscore');
var path = require('path');
var urljoin = require('url-join');
var humanize = require('humanize')

var FileItem = require('./FileItem.jsx')
var PathBreadcrumb = require('./PathBreadcrumb.jsx')
var UploadModal = require('./UploadModal.jsx');
var Icon = require('./Icon.jsx')
var Markdown = require('./Markdown.jsx')
var FilePreview = require('./FilePreview.jsx');


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
      historyDepth: 0,
      hidden: false,
      showUpload: false,
      readmeFile: null,
      pathname: decodeURI(location.pathname),
    }
  },
  loadFilesFromServer: function(){
    var that = this;
    $.ajax({
      url: that.state.pathname,
      data: {format: 'json'},
      dataType: 'json',
      success: function(data){
        data = _.sortBy(data, function(item){
          return item.type+':'+item.name;
        })
        this.setState({data: data})
        this.loadReadmeFromFiles(data)
      }.bind(this),
      error: function(xhr, status, err){
        console.log(status, err)
      }
    })
  },
  loadReadmeFromFiles: function(data){
    var readmes = _.filter(data, function(v){
      return _.contains(['readme', 'readme.md', 'readme.txt'], v.name.toLowerCase())
    });
    var readmeFile = _.max(readmes, function(v){
      var ext = path.extname(v.name.toLowerCase());
      if (ext == '.md') return 3;
      if (ext == '.txt') return 2;
      return 1;
    })
    if (readmeFile == -Infinity) {
      readmeFile = null;
    } else {
      readmeFile = path.join(this.state.pathname, readmeFile.name)
    }
    this.setState({readmeFile: readmeFile})
  },
  componentDidMount: function() {
    this.loadFilesFromServer();
    var that = this;
    window.onpopstate = function(event) {
      that.loadFilesFromServer();
    }
  },
  changePath: function(newPath, e){
    e.preventDefault()
    this.setState({pathname: newPath}, function(){
      this.loadFilesFromServer();
    })
    window.history.pushState({}, "", newPath);
  },
  render: function(){
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
                    <Button bsSize="xsmall" onClick={
                      ()=>this.setState({showUpload: true})
                    }>
                      Upload <Icon name="upload"/>
                    </Button>
                    <Button bsSize="xsmall"　onClick={
                      ()=>this.setState({hidden: !this.state.hidden})
                    }>
                      Show Hidden　{
                        this.state.hidden ? <i className="fa fa-eye"/> : <i className="fa fa-eye-slash"/>
                      }
                    </Button>
                  </ButtonToolbar>
                  
                  <UploadModal onHide={
                      ()=>this.setState({showUpload: false})}
                    onUpload={
                      ()=>this.loadFilesFromServer()
                    }
                    show={this.state.showUpload}/>
                </td>
              </tr>
              <tr>
                <th>
                  <Button bsSize="xsmall" 
                    href={path.dirname(this.state.pathname)}
                    onClick={(event)=>this.changePath(path.dirname(this.state.pathname), event)}>
                    <Icon name="arrow-up"/>
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
                { 
                  this.state.readmeFile ? (
                    <td colSpan={5}>
                      <FilePreview fileName={this.state.readmeFile}/>
                    </td>) : null
                }
              </tr>
            </tfoot>
          </Table>
        </Col>
      </Row>
    )
  }
})

module.exports = Explorer
