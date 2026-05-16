# ms2 — Message Queue (Pertemuan 4)

Microservice notifikasi menggunakan Node.js + Express.js + RabbitMQ. Publisher menerima request HTTP lalu mendistribusikan pesan ke tiga consumer (Email, SMS, FCM) melalui RabbitMQ fanout exchange.

**URL Produksi:** `http://mq.anla.works`

---

## Arsitektur

```
HTTP Client
    │
    ▼ POST /publish
┌─────────────┐      fanout exchange       ┌──────────────┐
│  Publisher  │ ──────────────────────────▶│  email.queue │──▶ Email Consumer
│ (Express.js)│                            ├──────────────┤
└─────────────┘                            │  sms.queue   │──▶ SMS Consumer
                                           ├──────────────┤
                                           │  fcm.queue   │──▶ FCM Consumer
                                           └──────────────┘
```

## Struktur Project

```
ms2.message-queue/
├── src/
│   ├── config/
│   │   └── rabbitmq.js       # Koneksi & setup exchange/queue
│   ├── publisher/
│   │   ├── publisher.js      # Logic publish ke RabbitMQ
│   │   └── routes.js         # Express route POST /publish
│   └── consumer/
│       ├── emailConsumer.js  # Consumer email.queue
│       ├── smsConsumer.js    # Consumer sms.queue
│       └── fcmConsumer.js    # Consumer fcm.queue
├── src/app.js                # Entry point (mode publisher/consumer/all)
├── Dockerfile
└── package.json
```

## Endpoints API

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/health` | Health check |
| `POST` | `/publish` | Kirim notifikasi ke semua queue |

### Contoh Request

```bash
curl -X POST http://mq.anla.works/publish \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "12345",
    "user_id": "67890",
    "content": "New order received",
    "timestamp": "2025-03-11T10:00:00Z"
  }'
```

**Respons:**
```json
{ "code": 200, "message": "Message published successfully" }
```

## Environment Variables

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `RABBITMQ_URL` | `amqp://guest:guest@localhost:5672` | URL koneksi RabbitMQ |
| `PORT` | `8081` | Port HTTP server |
| `APP_MODE` | `all` | Mode: `publisher`, `consumer`, atau `all` |

## Menjalankan Lokal

```bash
# Jalankan RabbitMQ
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:3-management-alpine

# Install dependencies
npm install

# Jalankan semua (publisher + consumer)
APP_MODE=all node src/app.js
```

## Deploy ke Kubernetes

```bash
kubectl apply -f k8s/ms2.message-queue/namespace.yaml
kubectl apply -f k8s/ms2.message-queue/configmap.yaml
kubectl apply -f k8s/ms2.message-queue/rabbitmq-deployment.yaml
kubectl apply -f k8s/ms2.message-queue/publisher-deployment.yaml
kubectl apply -f k8s/ms2.message-queue/consumer-deployment.yaml
kubectl apply -f k8s/ms2.message-queue/ingress.yaml
```

## Images yang Digunakan

| Image | Sumber | Fungsi |
|-------|--------|--------|
| `itsanla/mq-s6` | Docker Hub (custom) | Publisher & Consumer app |
| `rabbitmq:3-management-alpine` | Docker Hub (official) | Message broker |
