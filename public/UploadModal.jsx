var React = require('react');
var Dropzone = require('react-dropzone');
var Button = require('react-bootstrap').Button;
var Modal = require('react-bootstrap').Modal;
var request = require('superagent');
var Col = require('react-bootstrap').Col;
var Row = require('react-bootstrap').Row;


var UploadModal = React.createClass({
  getInitialState: function() {
    return {
      message: ""
    };
  },
  handleDrop: function(files){
    console.log("FILE", files)
    var that = this;
    var req = request.post(location.pathname);
    files.forEach((file)=> {
        req.attach('file', file, file.name);
    });
    req.end(function(err, res){
      if(res.ok) {
        console.log(res.body)
        that.setState({message: res.body.message})
      } else {
        that.setState({message: res.text});
      }
      var callback = that.props.onUpload;
      if (callback){
        callback(res);
      }
    })
    // req.end(callback);
  },
  render: function() {
    return (
      <div>        
        <Modal
          bsSize="small"
          show={this.props.show}
          onHide={this.props.onHide}
        >
          <Modal.Header closeButton>
            <Modal.Title className="text-center">File upload</Modal.Title>
            <Modal.Body>
              <div>
                <Row>
                  <Col md={6}>
                    <Dropzone onDrop={this.handleDrop}>
                      <div>Drop somethins here </div>
                    </Dropzone>
                  </Col>
                  <Col md={6} show={this.state.message != ""}>
                    {this.state.message}
                  </Col>
                </Row>
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button onClick={this.props.onHide}>Close</Button>
            </Modal.Footer>
          </Modal.Header>
        </Modal>
      </div>
    );
  }

});

module.exports = UploadModal;