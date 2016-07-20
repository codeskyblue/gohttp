var React = require('react');
var {Col, Row, ProgressBar, Modal, Button, Input} = require('react-bootstrap');
var Dropzone = require('react-dropzone');
var Icon = require('./Icon.jsx');
var request = require('superagent');

var Base64 = {_keyStr:"-+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789=",encode:function(e){var t="";var n,r,i,s,o,u,a;var f=0;e=Base64._utf8_encode(e);while(f<e.length){n=e.charCodeAt(f++);r=e.charCodeAt(f++);i=e.charCodeAt(f++);s=n>>2;o=(n&3)<<4|r>>4;u=(r&15)<<2|i>>6;a=i&63;if(isNaN(r)){u=a=64}else if(isNaN(i)){a=64}t=t+this._keyStr.charAt(s)+this._keyStr.charAt(o)+this._keyStr.charAt(u)+this._keyStr.charAt(a)}return t},decode:function(e){var t="";var n,r,i;var s,o,u,a;var f=0;e=e.replace(/[^A-Za-z0-9+/=]/g,"");while(f<e.length){s=this._keyStr.indexOf(e.charAt(f++));o=this._keyStr.indexOf(e.charAt(f++));u=this._keyStr.indexOf(e.charAt(f++));a=this._keyStr.indexOf(e.charAt(f++));n=s<<2|o>>4;r=(o&15)<<4|u>>2;i=(u&3)<<6|a;t=t+String.fromCharCode(n);if(u!=64){t=t+String.fromCharCode(r)}if(a!=64){t=t+String.fromCharCode(i)}}t=Base64._utf8_decode(t);return t},_utf8_encode:function(e){e=e.replace(/rn/g,"n");var t="";for(var n=0;n<e.length;n++){var r=e.charCodeAt(n);if(r<128){t+=String.fromCharCode(r)}else if(r>127&&r<2048){t+=String.fromCharCode(r>>6|192);t+=String.fromCharCode(r&63|128)}else{t+=String.fromCharCode(r>>12|224);t+=String.fromCharCode(r>>6&63|128);t+=String.fromCharCode(r&63|128)}}return t},_utf8_decode:function(e){var t="";var n=0;var r=c1=c2=0;while(n<e.length){r=e.charCodeAt(n);if(r<128){t+=String.fromCharCode(r);n++}else if(r>191&&r<224){c2=e.charCodeAt(n+1);t+=String.fromCharCode((r&31)<<6|c2&63);n+=2}else{c2=e.charCodeAt(n+1);c3=e.charCodeAt(n+2);t+=String.fromCharCode((r&15)<<12|(c2&63)<<6|c3&63);n+=3}}return t}}

var RDownloadModal = React.createClass({
  getInitialState : function () {
      return {
          url : "",
          message : "",
          percent : 0,
          downloading : false,
          timerid : 0,
          fsize : 0,
          fname : "",
      };
  },
  handleChange : function () {
      // This could also be done using ReactLink:
      // http://facebook.github.io/react/docs/two-way-binding-helpers.html
      this.setState({
          url : this.refs.input.getValue()
      });
  },
  sizeString : function (sz) {
      var us = "BKMGT";
      for (var i = 0; i < 5 && sz > 1000; i++) {
          sz = sz * 1.00 / 1024;
      }
      return this.formatString("{0} {1}", sz.toFixed(2), us[i]);
  },
  formatString: function() {
    if (arguments.length == 0)
        return null;
    var str = arguments[0];
    for ( var i = 1; i < arguments.length; i++) {
        var re = new RegExp('\\{' + (i - 1) + '\\}', 'gm');
        str = str.replace(re, arguments[i]);
    }
    return str;
  },
  update : function (flag) {
      if (flag && this.state.downloading) {
          var tid = setTimeout(this.onUpdate, 100);
          this.setState({
              timerid : tid
          });
      } else {
          clearTimeout(this.state.timerid);
          this.setState({
              timerid : 0
          });
      }
  },
  onUpdate : function () {
      var that = this;
      $.getJSON('/$wstat/' + this.state.fname, function (rs) {
          if (rs.downloaded == -1) {
              that.update(false);
              that.setState({
                  downloading : false,
                  message : "Failed to get downloading info"
              });
          } else {
              if (rs.downloaded >= that.state.fsize) {
                  that.setState({
                      percent : 100,
                      message : that.formatString("Download completed\n[ {0} ]\t{1}", that.sizeString(rs.downloaded), that.state.fname),
                      downloading: false
                  });
              } else {
                  that.setState({
                      percent : rs.downloaded * 100 / that.state.fsize,
                      message: that.formatString("[{0} / {2}] {1}", that.sizeString(rs.downloaded), that.state.fname, that.sizeString(that.state.fsize))
                  });
                  that.update(true);
              }
          }
      });
  },
  handleDownload : function () {
      console.log("download:", this.state.url);
      this.setState({
          downloading : true,
          percent : 0,
          fsize: 0,
          fname: "",
          message: "",
      });
      var that = this;
      var req = request.get('/$wget/' + Base64.encode(this.state.url));
      req.end(function (err, res) {
          if (res.ok) {
              var j = JSON.parse(res.text);
              console.log(j);
              that.setState({
                  fname : j.fname,
                  fsize : j.fsize
              });
              that.update(true);
          } else {
              that.setState({
                  message : res.text,
                  downloading : false
              });
          }
      })
  },
  onHide : function () {
      if (this.props.onHide) {
          this.props.onHide();
      }
      this.setState({
          percent : 0,
          message : "",
          url : "",
          downloading : false,
      });
      this.update(false);
  },
  render: function() {
    let downloading = this.state.downloading
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
                label="Arguments passed to wget:"
                ref="input"
                onChange={this.handleChange} />
              </form>
              <ProgressBar now={this.state.percent} label="%(percent)s%"/>
                {this.state.message ?
                  <div>
                    {this.state.message}
                  </div>: null
                }
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button bsStyle="primary" disabled={downloading} onClick={!downloading?this.handleDownload : null}>{downloading?'Downloading':'Download'}</Button>
              <Button onClick={this.onHide}>Close</Button>
            </Modal.Footer>
          </Modal.Header>
        </Modal>
      </div>
    );
  }

});

module.exports = RDownloadModal;