const { getChannel, QUEUES } = require('../config/rabbitmq');

function startEmailConsumer() {
  const channel = getChannel();

  channel.consume(QUEUES.EMAIL, (msg) => {
    if (!msg) return;

    const data = JSON.parse(msg.content.toString());
    console.log(`[Email Consumer] Menerima pesan:`);
    console.log(`  order_id : ${data.order_id}`);
    console.log(`  user_id  : ${data.user_id}`);
    console.log(`  content  : ${data.content}`);
    console.log(`  timestamp: ${data.timestamp}`);
    console.log(`  [SIMULASI] Email notifikasi terkirim ke user ${data.user_id}`);

    channel.ack(msg);
  }, { noAck: false });

  console.log(`[Email Consumer] Mendengarkan queue: ${QUEUES.EMAIL}`);
}

module.exports = { startEmailConsumer };
