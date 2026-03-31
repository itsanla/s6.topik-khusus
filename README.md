# Async Email Queue with Next.js + RabbitMQ

Project ini menggunakan arsitektur producer-consumer untuk pengiriman email asynchronous:

- Route handler `POST /api/email` bertindak sebagai producer (enqueue job email)
- Worker terpisah `workers/email-worker.ts` bertindak sebagai consumer (dequeue job + kirim email via SMTP)

Dengan pola ini, request HTTP bisa merespons cepat tanpa menunggu proses kirim email selesai.

## 1. Setup Environment Variable

Salin contoh env:

```bash
cp .env.example .env.local
```

Lalu isi nilai SMTP dan RabbitMQ sesuai environment Anda.

## 2. Install Dependencies

```bash
pnpm install
```

## 3. Jalankan Aplikasi dan Worker

Jalankan Next.js API producer:

```bash
pnpm dev
```

Jalankan worker email di terminal lain:

```bash
pnpm worker:email
```

## 4. Test Enqueue Email

Contoh request untuk memasukkan email ke queue:

```bash
curl -X POST http://localhost:3000/api/email \
	-H "Content-Type: application/json" \
	-d '{
		"to": "penerima@example.com",
		"subject": "Tes Async Email",
		"text": "Halo! Ini email dari antrian RabbitMQ."
	}'
```

Jika berhasil, API mengembalikan status `202` (Accepted), lalu worker akan memproses queue di background.

## 5. Struktur File Utama

- `app/api/email/route.ts`: endpoint enqueue email
- `lib/rabbitmq.ts`: koneksi RabbitMQ + helper publish/consume
- `lib/email.ts`: helper SMTP Nodemailer
- `workers/email-worker.ts`: consumer email queue
- `.env.example`: contoh konfigurasi env
