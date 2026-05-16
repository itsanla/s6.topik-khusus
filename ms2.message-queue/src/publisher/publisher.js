const { getChannel, EXCHANGE_NAME } = require('../config/rabbitmq');

function publish(message) {
  const channel = getChannel();
  const payload = Buffer.from(JSON.stringify(message));

  const ok = channel.publish(EXCHANGE_NAME, '', payload, {
    persistent    : true,
    contentType   : 'application/json',
    timestamp     : Math.floor(Date.now() / 1000),
  });

  if (!ok) throw new Error('Gagal publish ke exchange, channel penuh');
  return true;
}

module.exports = { publish };
