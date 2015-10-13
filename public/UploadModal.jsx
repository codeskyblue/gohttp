var React = require('react');
var Dropzone = require('react-dropzone');
var Button = require('react-bootstrap').Button;
var Modal = require('react-bootstrap').Modal;


var UploadModal = React.createClass({
  
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
              <div className="text-center">
                <Dropzone>
                  <div>Drop somethins here </div>
                </Dropzone>
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