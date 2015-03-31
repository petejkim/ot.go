(function () {
  "use strict";

  window.App = {
    conn: null
  }

  CodeMirror.fromTextArea(document.getElementById('code'), {
    lineNumbers: true,
    readOnly: true,
    mode: "go"
  });

  $('#join-btn').click(function(evt) {
    evt.preventDefault();
    $(this).attr({disabled: true});
    var $username = $('#join-form input[name=username]')
    $username.attr({disabled: true});
  });

  var url = [location.protocol.replace("http", "ws"), '//', location.host, '/ws'].join('')
  var conn = App.conn = new WebSocket(url);

  conn.onopen = function(evt) {
    $('#conn-status').text("Connected");
  };

  conn.onclose = function(evt) {
    $('#conn-status').text("Disconnected");
  };

  conn.onmessage = function(evt) {
    console.log(evt.data)
  };
}());
