# 🗃️ Database Schema Documentation

Dokumentasi ini menjelaskan struktur data yang digunakan dalam database PostgreSQL dan MongoDB untuk sistem blog.

---

## 🐘 PostgreSQL Schemas

### 1. `users` Table

Menyimpan informasi pengguna sistem.

| Column       | Type      | Constraints                               | Description               |
| ------------ | --------- | ----------------------------------------- | ------------------------- |
| `id`         | UUID      | Primary Key, Default: gen\_random\_uuid() | ID unik pengguna          |
| `name`       | String    | NOT NULL                                  | Nama lengkap pengguna     |
| `email`      | String    | NOT NULL, UNIQUE                          | Alamat email pengguna     |
| `username`   | String    | NOT NULL, UNIQUE                          | Nama pengguna (username)  |
| `password`   | String    | -                                         | Password terenkripsi      |
| `avatar`     | String    | -                                         | URL gambar profil         |
| `created_at` | Timestamp | Auto-generated                            | Waktu pembuatan akun      |
| `updated_at` | Timestamp | Auto-updated                              | Waktu terakhir diperbarui |

---

### 2. `metadata` Table

Menyimpan metadata artikel untuk keperluan moderasi dan tampilan publik.

| Column              | Type      | Constraints                               | Description                                 |
| ------------------- | --------- | ----------------------------------------- | ------------------------------------------- |
| `id`                | UUID      | Primary Key, Default: gen\_random\_uuid() | ID unik metadata                            |
| `article_id`        | String    | NOT NULL                                  | ID referensi ke koleksi artikel             |
| `title`             | String    | NOT NULL                                  | Judul artikel                               |
| `author`            | String    | -                                         | Nama penulis artikel                        |
| `moderation_status` | String    | -                                         | Status moderasi (pending/approved/rejected) |
| `created_at`        | Timestamp | Auto-generated                            | Waktu metadata dibuat                       |
| `updated_at`        | Timestamp | Auto-updated                              | Waktu metadata diubah                       |

---

## 🍃 MongoDB Schema

### 1. `articles` Collection

Menyimpan konten artikel dalam format tidak terstruktur.

| Field     | Type     | Description                               |
| --------- | -------- | ----------------------------------------- |
| `_id`     | ObjectID | ID unik dokumen MongoDB                   |
| `content` | String   | Konten artikel dalam format HTML/Markdown |
| `user_id` | String   | ID pengguna yang menulis artikel          |
| `status`  | String   | Status artikel (draft/published)          |

---

Silakan gunakan skema ini sebagai referensi saat melakukan migrasi data, validasi input, atau integrasi antar service.
