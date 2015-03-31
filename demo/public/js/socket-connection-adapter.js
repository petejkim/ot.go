(function () {
  'use strict';

  function SocketConnectionAdapter (conn) {
    this.conn = conn;

    var self = this;
    conn.on('quit', function (clientId) {
      self.trigger('client_left', clientId);
    });

    conn.on('join', function (data) {
      var clientId = data.client_id,
          name = data.username;
      self.trigger('set_name', clientId, name);
    });

    conn.on('ack', function () {
      self.trigger('ack');
    });

    conn.on('operation', function (data) {
      var clientId = data[0],
          operation = data[1],
          selection = data[2];
      self.trigger('operation', operation);
      self.trigger('selection', clientId, selection);
    });

    conn.on('selection', function (data) {
      var clientId = data[0],
          selection = data[1];
      self.trigger('selection', clientId, selection);
    });

    conn.on('reconnect', function () {
      self.trigger('reconnect');
    });
  }

  SocketConnectionAdapter.prototype.sendOperation = function (revision, operation, selection) {
    this.conn.send('operation', [revision, operation, selection]);
  };

  SocketConnectionAdapter.prototype.sendSelection = function (selection) {
    this.conn.send('selection', selection);
  };

  SocketConnectionAdapter.prototype.registerCallbacks = function (cb) {
    this.callbacks = cb;
  };

  SocketConnectionAdapter.prototype.trigger = function (event) {
    var args = Array.prototype.slice.call(arguments, 1);
    var action = this.callbacks && this.callbacks[event];
    if (action) { action.apply(this, args); }
  };

  window.ot.SocketConnectionAdapter = SocketConnectionAdapter;
}());
