const { getChannel, QUEUES } = require('../config/rabbitmq');

function startSmsConsumer() {
  const channel = getChannel();

  channel.consume(QUEUES.SMS, (msg) => {
    if (!msg) return;

    const data = JSON.parse(msg.content.toString());
    console.log(`[SMS Consumer] Menerima pesan:`);
    console.log(`  order_id : ${data.order_id}`);
    console.log(`  user_id  : ${data.user_id}`);
    console.log(`  content  : ${data.content}`);
    console.log(`  timestamp: ${data.timestamp}`);
    console.log(`  [SIMULASI] SMS notifikasi terkirim ke user ${data.user_id}`);

    channel.ack(msg);
  }, { noAck: false });

  console.log(`[SMS Consumer] Mendengarkan queue: ${QUEUES.SMS}`);
}

module.exports = { startSmsConsumer };
