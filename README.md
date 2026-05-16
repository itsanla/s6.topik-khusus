# S6 Topik Khusus — Microservices

Kumpulan microservice Go untuk mata kuliah Topik Khusus Semester 6. Setiap service berjalan independen dan di-orchestrate oleh Kubernetes (`k8s/*`).

## Struktur Monorepo

```
s6.topik-khusus/
├── k8s/                        # Kubernetes manifests semua service
│   └── track-method/
├── track-method/               # Pertemuan 3
├── message-queue/              # Pertemuan 4
├── nosql-mongodb/              # Pertemuan 5
├── langstats-indonesia/        # Pertemuan 6
├── elasticsearch-guide/        # Pertemuan 7
└── northwind-go/               # Pertemuan 8
```

## Daftar Service

| Pertemuan | Folder | Topik | Status |
|-----------|--------|-------|--------|
| 3 | [`track-method/`](./track-method) | Redis — Event Tracking dengan INCR | ✅ Selesai |
| 4 | [`message-queue/`](./message-queue) | Message Queue | 🔧 Belum |
| 5 | [`nosql-mongodb/`](./nosql-mongodb) | NoSQL MongoDB | 🔧 Belum |
| 6 | [`langstats-indonesia/`](./langstats-indonesia) | LangStats Indonesia | 🔧 Belum |
| 7 | [`elasticsearch-guide/`](./elasticsearch-guide) | Elasticsearch Guide | 🔧 Belum |
| 8 | [`northwind-go/`](./northwind-go) | Northwind Go — Clean Architecture API | 🔧 Belum |

## Stack Teknologi

- **Bahasa:** Go
- **Container:** Docker (Alpine)
- **Orchestration:** Kubernetes
- **Registry:** [Docker Hub — itsanla](https://hub.docker.com/u/itsanla)
- **CI/CD:** GitHub Actions

## Deploy Kubernetes

```bash
# Apply semua manifest sesuai service yang ingin di-deploy
kubectl apply -f k8s/<nama-service>/
```
