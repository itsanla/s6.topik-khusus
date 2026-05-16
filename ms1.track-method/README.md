# Track Method — Redis Microservice

Microservice berbasis Go untuk melacak event (click tracking, page view, dsb.) menggunakan Redis sebagai penyimpanan in-memory. Dibangun dengan Clean Architecture dan siap deploy ke Kubernetes.

**URL Produksi:** `http://track.anla.works`

---

## Arsitektur

```
track-method/
├── config/            # Konfigurasi aplikasi (env var)
├── domain/            # Interface & entity (kontrak)
├── repository/        # Implementasi akses Redis
├── usecase/           # Logika bisnis tracking
├── handler/           # HTTP handler (Gin)
├── middleware/        # Logger & CORS middleware
├── k8s/               # Kubernetes manifests
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── redis-deployment.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── hpa.yaml
├── main.go
├── Dockerfile
└── go.mod
```

## Cara Kerja

Setiap kali endpoint `POST /api/v1/track` dipanggil dengan nama event, service akan menjalankan perintah `INCR` ke Redis dengan key `track:event:{nama_event}`. Counter bertambah atomik dan thread-safe secara native di Redis.

## Endpoints API

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/health` | Health check service dan Redis |
| `POST` | `/api/v1/track` | Catat / increment sebuah event |
| `GET` | `/api/v1/track` | Ambil statistik semua event |
| `GET` | `/api/v1/track/:event` | Ambil statistik satu event |
| `DELETE` | `/api/v1/track/:event` | Reset counter event |

### Contoh Request

**Catat event:**
```bash
curl -X POST http://track.anla.works/api/v1/track \
  -H "Content-Type: application/json" \
  -d '{"event": "click_button_login"}'
```

**Respons:**
```json
{
  "success": true,
  "message": "event berhasil dicatat",
  "data": {
    "event_name": "click_button_login",
    "count": 3
  }
}
```

**Ambil semua statistik:**
```bash
curl http://track.anla.works/api/v1/track
```

## Environment Variables

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `REDIS_ADDR` | `localhost:6379` | Alamat Redis |
| `REDIS_PASSWORD` | `` | Password Redis |
| `REDIS_DB` | `0` | Database Redis |
| `PORT` | `8080` | Port HTTP server |
| `APP_ENV` | `development` | Mode aplikasi |

## Menjalankan Lokal

```bash
# Jalankan Redis
docker run -d -p 6379:6379 redis:7-alpine

# Build dan jalankan
docker build -t track-method .
docker run -p 8080:8080 \
  -e REDIS_ADDR=host.docker.internal:6379 \
  track-method
```

## Deploy ke Kubernetes

> Manifest Kubernetes berada di root project: `k8s/track-method/`

```bash
# Apply semua manifest secara berurutan (dari root project)
kubectl apply -f k8s/track-method/namespace.yaml
kubectl apply -f k8s/track-method/configmap.yaml
kubectl apply -f k8s/track-method/secret.yaml
kubectl apply -f k8s/track-method/redis-deployment.yaml
kubectl apply -f k8s/track-method/deployment.yaml
kubectl apply -f k8s/track-method/service.yaml
kubectl apply -f k8s/track-method/ingress.yaml
kubectl apply -f k8s/track-method/hpa.yaml

# Cek status pod
kubectl get pods -n track-method

# Cek logs
kubectl logs -l app=track-method -n track-method
```

## Prompt AI yang Digunakan

> "Buatkan microservice Golang menggunakan Redis sebagai in-memory store dengan metode INCR untuk tracking event. Gunakan Clean Architecture (domain, repository, usecase, handler). Tambahkan HTTP REST API menggunakan Gin framework, konfigurasi berbasis environment variable, Dockerfile multi-stage build Alpine, dan Kubernetes manifests (Deployment 3 replika, Service, Ingress, HPA). Domain yang digunakan adalah track.anla.works."
