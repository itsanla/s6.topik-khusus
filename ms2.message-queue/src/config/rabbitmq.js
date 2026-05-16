const amqp = require('amqplib');

const RABBITMQ_URL = process.env.RABBITMQ_URL || 'amqp://guest:guest@localhost:5672';

const EXCHANGE_NAME = 'notifications';
const EXCHANGE_TYPE  = 'fanout';

const QUEUES = {
  EMAIL : 'email.queue',
  SMS   : 'sms.queue',
  FCM   : 'fcm.queue',
};

let connection = null;
let channel    = null;

async function connect(retries = 10, delay = 3000) {
  for (let i = 1; i <= retries; i++) {
    try {
      connection = await amqp.connect(RABBITMQ_URL);
      channel    = await connection.createChannel();

      await channel.assertExchange(EXCHANGE_NAME, EXCHANGE_TYPE, { durable: true });

      for (const q of Object.values(QUEUES)) {
        await channel.assertQueue(q, { durable: true });
        await channel.bindQueue(q, EXCHANGE_NAME, '');
      }

      console.log(`[RabbitMQ] Terhubung ke ${RABBITMQ_URL}`);

      connection.on('error', (err) => {
        console.error('[RabbitMQ] Koneksi error:', err.message);
      });
      connection.on('close', () => {
        console.warn('[RabbitMQ] Koneksi tertutup, mencoba reconnect...');
        setTimeout(() => connect(), 5000);
      });

      return channel;
    } catch (err) {
      console.warn(`[RabbitMQ] Percobaan ${i}/${retries} gagal: ${err.message}`);
      if (i === retries) throw err;
      await new Promise((r) => setTimeout(r, delay));
    }
  }
}

function getChannel() {
  if (!channel) throw new Error('RabbitMQ belum terhubung');
  return channel;
}

module.exports = { connect, getChannel, EXCHANGE_NAME, QUEUES };
