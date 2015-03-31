(function () {
  'use strict';

  function SocketConnection (url) {
    this.url = url;
    this.ws = null;

    var ws = this.ws = new WebSocket(this.url);
    var self = this;

    ws.onopen = function (evt) {
      self.emit('open', evt);
    };

    ws.onclose = function (evt) {
      self.emit('close', evt);
    };

    ws.onmessage = function (evt) {
      var m = JSON.parse(evt.data);

      if (m && m.event) {
        self.emit(m.event, m.data, evt);
      }
    };
  }

  SocketConnection.prototype = new EventEmitter;

  SocketConnection.prototype.send = function(eventName, data) {
    this.ws.send(JSON.stringify({event: eventName, data: data}));
  };

  window.SocketConnection = SocketConnection;
}());
