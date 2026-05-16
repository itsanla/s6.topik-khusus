const { getChannel, QUEUES } = require('../config/rabbitmq');

function startFcmConsumer() {
  const channel = getChannel();

  channel.consume(QUEUES.FCM, (msg) => {
    if (!msg) return;

    const data = JSON.parse(msg.content.toString());
    console.log(`[FCM Consumer] Menerima pesan:`);
    console.log(`  order_id : ${data.order_id}`);
    console.log(`  user_id  : ${data.user_id}`);
    console.log(`  content  : ${data.content}`);
    console.log(`  timestamp: ${data.timestamp}`);
    console.log(`  [SIMULASI] Push notification FCM terkirim ke user ${data.user_id}`);

    channel.ack(msg);
  }, { noAck: false });

  console.log(`[FCM Consumer] Mendengarkan queue: ${QUEUES.FCM}`);
}

module.exports = { startFcmConsumer };
