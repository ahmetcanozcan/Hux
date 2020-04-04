var Hux = (() => {
  const sep = "]];|;[[";
  let ws;
  let events = {};

  function constructor(URL = "ws://localhost:8080/ws/hux") {
    events = {};
    ws = new WebSocket(URL);
    const reader = new FileReader();
    var self = this;
    ws.onmessage = function (msg) {
      var listener = e => {
        const text = e.srcElement.result;
        var [name, data] = text.split(sep);
        if (!name || !data) {
          return console.log('Invalid message.');
        }
        // Invoke event handler.
        var f = events[name];
        if (f) {
          f(data);
        }
        // After All clear event.
        reader.removeEventListener('loadend', listener)
      }
      reader.addEventListener('loadend', listener);
      reader.readAsText(msg.data);

    }
    ws.onclose = function () {
      var f = events['close'];
      if (f) {
        f();
      }

    }
    ws.onopen = function () {
      var f = events['open'];
      if (f) {
        f();
      }
    }
  }
  /**
   * @param {string} name event name
   * @param {Function} handler
   */
  constructor.prototype.on = function (name, handler) {
    events[name] = handler;
  }



  constructor.prototype.emit = function (name, data) {
    const msg = `${name}${sep}${data}`;
    ws.send(msg)
  }

  return constructor
})();
