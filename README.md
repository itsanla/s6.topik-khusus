# S6 Topik Khusus — Microservices

Kumpulan microservice Go untuk mata kuliah Topik Khusus Semester 6. Setiap service berjalan independen dan di-orchestrate oleh Kubernetes (`k8s/*`).

## Struktur Monorepo

```
s6.topik-khusus/
├── k8s/                            # Kubernetes manifests semua service
│   └── ms1.track-method/
├── ms1.track-method/               # Pertemuan 3 — Redis
├── ms2.message-queue/              # Pertemuan 4 — Message Queue
├── ms3.nosql-mongodb/              # Pertemuan 5 — NoSQL MongoDB
├── ms4.langstats-indonesia/        # Pertemuan 6 — LangStats Indonesia
├── ms5.elasticsearch-guide/        # Pertemuan 7 — Elasticsearch Guide
└── ms6.northwind-go/               # Pertemuan 8 — Northwind Go
```

## Daftar Service

| # | Folder | Pertemuan | Topik | Status |
|---|--------|-----------|-------|--------|
| ms1 | [`ms1.track-method/`](./ms1.track-method) | 3 | Redis — Event Tracking dengan INCR | ✅ Selesai |
| ms2 | [`ms2.message-queue/`](./ms2.message-queue) | 4 | Message Queue | 🔧 Belum |
| ms3 | [`ms3.nosql-mongodb/`](./ms3.nosql-mongodb) | 5 | NoSQL MongoDB | 🔧 Belum |
| ms4 | [`ms4.langstats-indonesia/`](./ms4.langstats-indonesia) | 6 | LangStats Indonesia | 🔧 Belum |
| ms5 | [`ms5.elasticsearch-guide/`](./ms5.elasticsearch-guide) | 7 | Elasticsearch Guide | 🔧 Belum |
| ms6 | [`ms6.northwind-go/`](./ms6.northwind-go) | 8 | Northwind Go — Clean Architecture API | 🔧 Belum |

## Stack Teknologi

- **Bahasa:** Go
- **Container:** Docker (Alpine)
- **Orchestration:** Kubernetes
- **Registry:** [Docker Hub — itsanla](https://hub.docker.com/u/itsanla)
- **CI/CD:** GitHub Actions

## Deploy Kubernetes

```bash
# Apply semua manifest sesuai service yang ingin di-deploy
kubectl apply -f k8s/ms1.track-method/
```
