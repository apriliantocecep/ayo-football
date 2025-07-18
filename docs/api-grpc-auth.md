# 📄 AuthService - gRPC API Documentation

## Overview

`AuthService` adalah gRPC service untuk autentikasi dan manajemen pengguna. Service ini terdiri dari empat method utama:

* `Login`: Autentikasi pengguna.
* `Register`: Registrasi pengguna baru.
* `ValidateToken`: Validasi JWT token.
* `GetUserById`: Mendapatkan data pengguna berdasarkan ID.

---

## 🛠️ Service: `AuthService`

### 1. `Login`

Autentikasi pengguna menggunakan identifier (email/username) dan password.

* **Request**: [`LoginRequest`](#loginrequest)
* **Response**: [`LoginResponse`](#loginresponse)

### 2. `Register`

Registrasi pengguna baru.

* **Request**: [`RegisterRequest`](#registerrequest)
* **Response**: [`RegisterResponse`](#registerresponse)

### 3. `ValidateToken`

Validasi token JWT dan mengembalikan ID pengguna.

* **Request**: [`ValidateTokenRequest`](#validatetokenrequest)
* **Response**: [`ValidateTokenResponse`](#validatetokenresponse)

### 4. `GetUserById`

Mengambil informasi pengguna berdasarkan ID.

* **Request**: [`GetUserByIdRequest`](#getuserbyidrequest)
* **Response**: [`User`](#user)

---

## 📦 Messages - AuthService

### LoginRequest

| Field        | Type   | Description                  |
| ------------ | ------ | ---------------------------- |
| `identifier` | string | Email atau username pengguna |
| `password`   | string | Password pengguna            |

### LoginResponse

| Field          | Type   | Description                     |
| -------------- | ------ | ------------------------------- |
| `access_token` | string | Token JWT yang valid untuk sesi |
| `expires_at`   | string | Waktu kadaluarsa token          |

### RegisterRequest

| Field      | Type   | Description           |
| ---------- | ------ | --------------------- |
| `email`    | string | Email pengguna        |
| `password` | string | Password pengguna     |
| `name`     | string | Nama lengkap pengguna |

### RegisterResponse

| Field      | Type   | Description              |
| ---------- | ------ | ------------------------ |
| `user_id`  | string | ID unik pengguna         |
| `username` | string | Username yang dihasilkan |

### ValidateTokenRequest

| Field   | Type   | Description               |
| ------- | ------ | ------------------------- |
| `token` | string | JWT token yang divalidasi |

### ValidateTokenResponse

| Field     | Type   | Description            |
| --------- | ------ | ---------------------- |
| `user_id` | string | ID pengguna yang valid |

### GetUserByIdRequest

| Field | Type   | Description        |
| ----- | ------ | ------------------ |
| `id`  | string | ID pengguna target |

### User

| Field        | Type   | Description          |
| ------------ | ------ | -------------------- |
| `id`         | string | ID pengguna          |
| `name`       | string | Nama pengguna        |
| `email`      | string | Email pengguna       |
| `username`   | string | Username pengguna    |
| `created_at` | string | Waktu dibuatnya akun |
