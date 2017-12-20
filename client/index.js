const { Client } = require('faye-websocket')

const client = new Client('ws://localhost:3000')

const send = (data) => {
  client.send(JSON.stringify(data))
}

client.onopen = function() {
  console.log('[open]', client.headers);
  let i = 0
  setInterval(() => send({ event: 'test', body: ++i }), 1000)
};

client.onclose = function(close) {
  console.log('[close]', close.code, close.reason);
};

client.onerror = function(error) {
  console.log('[error]', error.message);
};

client.onmessage = function(message) {
  console.log('[message]', message.data);
};
