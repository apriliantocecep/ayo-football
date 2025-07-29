# Setup Guide - Content Publisher Microservice

Dokumen ini menjelaskan cara menjalankan sistem backend Content Publisher berbasis microservices secara lokal menggunakan Docker Compose.

---

## 🔧 Requirements

Sebelum menjalankan proyek ini, pastikan kamu sudah menginstal:

- [Docker](https://www.docker.com/)
- [Go](https://golang.org/dl/) (jika ingin build manual)
- [protoc](https://grpc.io/docs/languages/go/quickstart/) (quick start untuk golang jika ingin generate kode dari `.proto`)

---

# Step 1 - Clone Repository

```shell
git clone https://github.com/apriliantocecep/ayo-football.git
```

# Step 2 -  Jalankan docker compose utama

Pada langkah ini, kamu harus menjalankan perintah docker compose untuk membuat provisioning microservice.

Termasuk didalamnya yaitu RabbitMQ, Database Postgre dan MongoDb.

Jalankan perintah berikut:

```shell
docker compose up -d --build
```

Catatan:

Jika terdapat konfigurasi tambahan atau modifikasi pada file docker-compose.yaml, lebih baik jalankan perintah berikut terlebih dahulu

```shell
docker compose down -v
```

# Step 3 -  Management

- Buka http://localhost:8200/ untuk menjalankan Vault. Gunakan token `root` untuk masuk.
- Buka http://localhost:15672/ untuk menjalankan RabbitMQ Web Management. Gunakan username dan password `guest`
- Buka http://localhost:3000/ untuk menjalankan Grafana Dashboard. Buka `Tempo` di menu Explorer
- Buka http://localhost:8500/ untuk menjalankan Consul (Service Discovery)
- Buka http://localhost:8080/ untuk menjalankan Traefik (Load Balancer)

# Step 4 - REST API Documentation

| Item                     | Keterangan                                                                                                           |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------- |
| **Postman Collection**   | Import file: [`Posfin Blog REST API.postman_collection.json`](../Posfin%20Blog%20REST%20API.postman_collection.json) |
| **Default API Base URL** | `http://localhost:8000`                                                                                              |
| **Postman Environment**  | - `url`: `http://localhost:8000`<br>- `token`: *auto-generate saat hit endpoint `/login`*                            |
