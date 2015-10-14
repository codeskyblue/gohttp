var React = require('react');
var Dropzone = require('react-dropzone');
var Button = require('react-bootstrap').Button;
var Modal = require('react-bootstrap').Modal;
var request = require('superagent');
var Col = require('react-bootstrap').Col;
var Row = require('react-bootstrap').Row;
var ProgressBar = require('react-bootstrap').ProgressBar;
var Icon = require('./Icon.jsx');

var UploadModal = React.createClass({
  getInitialState: function() {
    return {
      message: "",
      percent: 0,
    };
  },
  handleDrop: function(files){
    console.log("FILE", files)
    var that = this;
    var req = request.post(location.pathname);
    files.forEach((file)=> {
        req.attach('file', file, file.name);
    });
    req
      .on('progress', function(e){
        that.setState({percent: e.percent})
      })
      .end(function(err, res){
        if(res.ok) {
          that.setState({message: res.body.message})
        } else {
          that.setState({message: res.text});
        }
        var callback = that.props.onUpload;
        if (callback){
          callback(res);
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
            <Modal.Title className="text-center">File upload</Modal.Title>
            <Modal.Body>
              <div>
                <Dropzone onDrop={this.handleDrop} className='dropzone' activeClassName='dropzone-active'>
                  Drop or click to upload
                </Dropzone>
              
                
                <h4>Progress</h4>
                <ProgressBar now={this.state.percent} label="%(percent)s%"/>
                {this.state.message ?
                  <div>
                    {this.state.message}
                  </div>: null
                }
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button onClick={this.onHide}>Close</Button>
            </Modal.Footer>
          </Modal.Header>
        </Modal>
      </div>
    );
  }

});

module.exports = UploadModal;