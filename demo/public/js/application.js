(function () {
  'use strict';

  window.App = {
    conn: null,
    cm: null
  };

  App.cm = CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    readOnly: 'nocursor',
    mode: 'go'
  });

  $('#join-btn').click(function (evt) {
    evt.preventDefault();
    $(this).attr({disabled: true});
    var $username = $('#join-form input[name=username]');
    $username.attr({disabled: true});
    App.conn.send('join', { username: $username.val() });
  });

  var url = [location.protocol.replace('http', 'ws'), '//', location.host, '/ws'].join('');
  var conn = App.conn = new SocketConnection(url);

  conn.on('open', function () {
    $('#conn-status').text('Connected');
    $('#join-btn').attr({ disabled: false});
  });

  conn.on('close', function (evt) {
    $('#conn-status').text('Disconnected');
  });

  conn.on('doc', function(data) {
    App.cm.setValue(data.document);
    var serverAdapter = new ot.SocketConnectionAdapter(conn);
    var editorAdapter = new ot.CodeMirrorAdapter(App.cm);
    App.client = new ot.EditorClient(data.revision, data.clients, serverAdapter, editorAdapter);
  });

  conn.on('registered', function(clientId) {
    App.cm.setOption('readOnly', false);
  });

  conn.on('join', function(data) {
    console.log(data);
  });

  conn.on('quit', function(data) {
    console.log(data);
  });
}());
