# Architecture Overview

   Sistem ini adalah backend berbasis microservices untuk platform Content Publisher. Sistem memungkinkan user mengirimkan artikel (dalam format WYSIWYG HTML) dan menangani penyimpanan metadata, moderasi konten, serta pemrosesan concurrent menggunakan gRPC dan RabbitMQ.

![](assets/architecture.png)

   ## Service

   - **Gateway Service**: Menerima permintaan dari klien (misalnya frontend) dan meneruskannya via gRPC ke `service yang dibutuhkan`.
   - **Auth Service**: Menangani autentikasi user
   - **Article Service**: Menangani logika utama:
       - Menyimpan metadata (judul, penulis, waktu) ke PostgreSQL dengan GORM
       - Menyimpan konten HTML ke MongoDB
       - Mengirim tugas moderasi ke RabbitMQ
   - **Moderation Service**: Mendengarkan antrean RabbitMQ, menyimulasikan proses moderasi konten, dan (opsional) memperbarui status konten atau mencatat hasil.

   ## Teknologi

   - **gRPC** untuk komunikasi antar gateway dan service lain
   - **PostgreSQL** untuk penyimpanan metadata terstruktur
   - **MongoDB** untuk menyimpan konten artikel (WYSIWYG HTML)
   - **RabbitMQ** sebagai message broker untuk pipeline moderasi asinkron
   - **Docker Compose** untuk menjalankan semua layanan dan dependensi
   - **Vault** untuk menjalankan secret management

   ## Deployment

   Semua layanan dan dependensinya dijalankan dalam container Docker yang dikoordinasi oleh `docker-compose`. Konfigurasi environment menggunakan file `.env`.
   