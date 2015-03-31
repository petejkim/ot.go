(function () {
  'use strict';

  window.App = {
    conn: null,
    cm: null,
    users: [],

    updateUsers: function () {
      $('#users').text(this.users.join(', '));
    }
  };

  App.cm = CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    readOnly: true,
    mode: 'go'
  });

  $('#join-btn').click(function (evt) {
    evt.preventDefault();
    $(this).attr({disabled: true});
    var $username = $('#join-form input[name=username]');
    $username.attr({disabled: true});
    App.conn.send(JSON.stringify({event: 'join', data: { username: $username.val() }}));
  });

  var url = [location.protocol.replace('http', 'ws'), '//', location.host, '/ws'].join('');
  var conn = App.conn = new SocketConnection(url);

  conn.on('open', function () {
    $('#conn-status').text('Connected');
    $('#join-btn').attr({ disabled: false});
  });

  conn.on('close', function () {
    $('#conn-status').text('Disconnected');
  });

  conn.on('doc', function (data) {
    App.cm.setValue(data.document);
    if (data.clients) {
      App.users = App.users.concat(data.clients);
      App.updateUsers();
    }
  });

  conn.on('join', function (data) {
    var clientId = data.client_id;
    if (clientId) {
      App.users.push(clientId);
      App.updateUsers();
    }
  });

  conn.on('quit', function (data) {
    var clientId = data.client_id;
    if (clientId) {
      var i = App.users.indexOf(clientId);
      if (i !== -1) {
        App.users.splice(i, 1);
      }
      App.updateUsers();
    }
  });
}());
