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
  var conn = App.conn = new WebSocket(url);

  conn.onopen = function (evt) {
    $('#conn-status').text('Connected');
    $('#join-btn').attr({ disabled: false});
  };

  conn.onclose = function (evt) {
    $('#conn-status').text('Disconnected');
  };

  conn.onmessage = function (evt) {
    var m = JSON.parse(evt.data);
    console.log(m);

    switch (m.event) {
    case 'doc':
      App.cm.setValue(m.data.document);
      if (m.data.clients) {
        App.users = App.users.concat(m.data.clients);
        App.updateUsers();
      }
      break;
    case 'join':
      var clientId = m.data.client_id;
      if (clientId) {
        App.users.push(clientId);
        App.updateUsers();
      }
      break;
    case 'quit':
      var clientId = m.data.client_id;
      if (clientId) {
        var i = App.users.indexOf(clientId);
        if (i !== -1) {
          App.users.splice(i, 1);
        }
        App.updateUsers();
      }
      break;
    }
  };
}());