var React = require('react');
var {Col, Row, ProgressBar, Modal, Button, Input} = require('react-bootstrap');
var Dropzone = require('react-dropzone');
var Icon = require('./Icon.jsx');
var request = require('superagent');


var RDownloadModal = React.createClass({
  getInitialState: function() {
    return {
      url: "",
      message: "",
      percent: 0,
    };
  },
  handleChange: function() {
    // This could also be done using ReactLink:
    // http://facebook.github.io/react/docs/two-way-binding-helpers.html
    this.setState({
      url: this.refs.input.getValue()
    });
  },
  handleDownload: function() {
    console.log("download:", this.state.url)
    var that = this;
    var req = request.get(location.protocol+'//'+location.host+'/$wget/'+this.state.url);
    req.end(function(err, res){
        if(res.ok) {
          that.onHide();
        } else {
          that.setState({message: res.text});
        }
      })
  },
  onHide: function(){
    if(this.props.onHide){
      this.props.onHide();
    }
    this.setState({
      percent: 0,
      message: "",
      url: "",
    })
  },
  render: function() {
    return (
      <div>        
        <Modal
          bsSize="small"
          show={this.props.show}
          onHide={this.onHide}
        >
          <Modal.Header closeButton>
            <Modal.Title className="text-center">Remote download using wget</Modal.Title>
            <Modal.Body>
              <div>
              <form>
                <Input
                type="text"
                value={this.state.url}
                placeholder="http:// or https://"
                label="URL:"
                ref="input"
                onChange={this.handleChange} />
              </form>

                {this.state.message ?
                  <div>
                    {this.state.message}
                  </div>: null
                }
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button bsStyle="primary" onClick={this.handleDownload}>Download</Button>
              <Button onClick={this.onHide}>Close</Button>
            </Modal.Footer>
          </Modal.Header>
        </Modal>
      </div>
    );
  }

});

module.exports = RDownloadModal;