const express = require('express');
const { publish } = require('./publisher');

const router = express.Router();

router.post('/publish', (req, res) => {
  const { order_id, user_id, content, timestamp } = req.body;

  if (!order_id || !user_id || !content) {
    return res.status(400).json({
      code    : 400,
      message : 'Field order_id, user_id, dan content wajib diisi',
    });
  }

  const message = {
    order_id,
    user_id,
    content,
    timestamp: timestamp || new Date().toISOString(),
  };

  try {
    publish(message);
    console.log(`[Publisher] Pesan dikirim ke exchange:`, message);
    res.json({ code: 200, message: 'Message published successfully' });
  } catch (err) {
    console.error('[Publisher] Gagal publish:', err.message);
    res.status(500).json({ code: 500, message: err.message });
  }
});

module.exports = router;
