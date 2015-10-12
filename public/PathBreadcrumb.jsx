var React = require('react');
var Breadcrumb = require('react-bootstrap').Breadcrumb;
var BreadcrumbItem = require('react-bootstrap').BreadcrumbItem;
var _ = require('underscore');
var urljoin = require('url-join')

var PathBreadcrumb = React.createClass({
  render: function(){
    var paths = (this.props.data || location.pathname).split('/');
    paths = _.filter(paths, function(v){
    	return v !== "";
    })
    var itemPaths = [];
    var currPath = '/';
    paths.forEach(function(name){
      var newPath = urljoin(currPath, name);
      currPath = newPath;
      itemPaths.push({
        directory: currPath,
        name: name
      })
    })
    var that = this;
    var items = itemPaths.map(function(subPath){
      return (
        <BreadcrumbItem 
        	key={subPath.directory} 
        	href={subPath.directory} 
        	onClick={(event)=>that.props.onClick(subPath.directory, event)}>
          {subPath.name}
        </BreadcrumbItem>
      )
    })
    var homeLink = "/"
    return (
      <Breadcrumb>
        <BreadcrumbItem onClick={(event)=>this.props.onClick("/", event)} href="/">
          <i className="fa fa-home"/>
        </BreadcrumbItem>
        {items}
      </Breadcrumb>
    )
  }
});

module.exports = PathBreadcrumb;