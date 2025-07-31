# ⚽︎ Ayo Football - Backend Microservices

Sistem backend ini dibangun menggunakan arsitektur microservices untuk aplikasi Football Manager

![](docs/assets/techstack.png)

# 🔌 Teknologi
- Go (Golang)
- gRPC
- PostgreSQL
- MongoDB
- RabbitMQ
- Docker & Docker Compose
- Vault

# 🤖Exposed Ports

| Service Name      | Container Name    | Port (Host\:Container) | Keterangan                                                   |
|-------------------|-------------------|------------------------|--------------------------------------------------------------|
| `vault`           | `vault-dev`       | `8200:8200`            | Vault dev UI dan API                                         |
| `rabbitmq-server` | `rabbitmq-server` | `5672:5672`            | AMQP (RabbitMQ message broker)                               |
|                   |                   | `15672:15672`          | RabbitMQ management UI                                       |
| `postgres-server` | `postgres-server` | `5432:5432`            | PostgreSQL default port                                      |
| `mongodb-server`  | `mongodb`         | `27017:27017`          | MongoDB default port                                         |
| `vault-init`      | `vault-init`      | *(tidak exposed)*      | Container init satu kali jalan (init script)                 |
| `gateway-srv`     | `auth-srv`        | `8000:8000`            | External-only, expose port ke host sebagai base url REST API |
| `auth-srv`        | `auth-srv`        | *(tidak exposed):8001* | Internal-only, tidak expose port ke host                     |
| `player-srv`      | `player-srv`      | *(tidak exposed):8002* | Internal-only, tidak expose port ke host                     |
| `match-srv`       | `match-srv`       | *(tidak exposed):8003* | Internal-only, tidak expose port ke host                     |
| `team-srv`        | `team-srv`        | *(tidak exposed):8004* | Internal-only, tidak expose port ke host                     |
# ⚙️ Setup Guide

Panduan ini berisi langkah demi langkah untuk menjalankan sistem Ayo Football secara lokal menggunakan Docker. Di dalamnya terdapat instruksi mulai dari cloning repository, konfigurasi environment, hingga menjalankan seluruh service.

👉 Silakan baca file [setup-guide.md](./docs/setup-guide.md) untuk informasi lengkap.

# 🏗️ Architecture

Ayo Football menggunakan pendekatan microservices architecture yang terdiri dari beberapa service terpisah, seperti authentication, player, match, dan lainnya. Setiap service berkomunikasi melalui gRPC, dan dikelola melalui API Gateway berbasis REST.

Struktur arsitektur ini dirancang untuk skalabilitas, modularitas, dan kemudahan pengembangan.