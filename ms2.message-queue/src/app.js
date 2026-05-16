const express = require('express');
const { connect } = require('./config/rabbitmq');
const publisherRoutes = require('./publisher/routes');
const { startEmailConsumer } = require('./consumer/emailConsumer');
const { startSmsConsumer }   = require('./consumer/smsConsumer');
const { startFcmConsumer }   = require('./consumer/fcmConsumer');

const PORT     = process.env.PORT     || 8081;
const APP_MODE = process.env.APP_MODE || 'all';

const app = express();
app.use(express.json());

app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'mq-s6', mode: APP_MODE });
});

async function main() {
  await connect();

  if (APP_MODE === 'publisher' || APP_MODE === 'all') {
    app.use('/', publisherRoutes);
    app.listen(PORT, () => {
      console.log(`[Publisher] HTTP server berjalan di port ${PORT}`);
    });
  }

  if (APP_MODE === 'consumer' || APP_MODE === 'all') {
    startEmailConsumer();
    startSmsConsumer();
    startFcmConsumer();
    console.log('[Consumer] Semua consumer aktif dan mendengarkan queue');
  }
}

main().catch((err) => {
  console.error('[App] Gagal start:', err.message);
  process.exit(1);
});

process.on('SIGINT',  () => { console.log('Shutting down...'); process.exit(0); });
process.on('SIGTERM', () => { console.log('Shutting down...'); process.exit(0); });
