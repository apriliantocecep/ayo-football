# 📄 ArticleService - gRPC API Documentation

## Overview

`ArticleService` adalah gRPC service untuk mengelola artikel dalam sistem blogging. Service ini terdiri dari tiga method utama:

* `SubmitArticle`: Mengirim artikel baru.
* `PublishArticle`: Mempublikasikan artikel.
* `GetArticles`: Mengambil daftar artikel milik pengguna.

---

## 🛠️ Service: `ArticleService`

### 1. `SubmitArticle`

Mengirim artikel baru yang ditulis oleh pengguna.

* **Request**: [`SubmitArticleRequest`](#submitarticlerequest)
* **Response**: [`SubmitArticleResponse`](#submitarticleresponse)

### 2. `PublishArticle`

Mempublikasikan artikel yang sebelumnya telah dikirim.

* **Request**: [`PublishArticleRequest`](#publisharticlerequest)
* **Response**: [`PublishArticleResponse`](#publisharticleresponse)

### 3. `GetArticles`

Mengambil semua artikel milik pengguna berdasarkan ID pengguna.

* **Request**: [`GetArticlesRequest`](#getarticlesrequest)
* **Response**: [`GetArticlesResponse`](#getarticlesresponse)

---

## 📦 Messages

### SubmitArticleRequest

Permintaan untuk mengirim artikel baru.

| Field          | Type   | Description                  |
| -------------- | ------ | ---------------------------- |
| `title`        | string | Judul artikel                |
| `author`       | string | Nama penulis                 |
| `html_content` | string | Konten HTML dari artikel     |
| `user_id`      | string | ID pengguna yang mengirimkan |

---

### SubmitArticleResponse

Respons setelah artikel dikirim.

| Field        | Type   | Description                                  |
| ------------ | ------ | -------------------------------------------- |
| `article_id` | string | ID unik dari artikel                         |
| `status`     | string | Status pengiriman (`success`, `failed`, dll) |

---

### PublishArticleRequest

Permintaan untuk mempublikasikan artikel.

| Field        | Type   | Description                    |
| ------------ | ------ | ------------------------------ |
| `article_id` | string | ID artikel yang akan dipublish |
| `user_id`    | string | ID pengguna                    |

---

### PublishArticleResponse

Respons setelah mencoba mempublikasikan artikel.

| Field    | Type   | Description                                   |
| -------- | ------ | --------------------------------------------- |
| `status` | string | Status publikasi (`published`, `failed`, dll) |

---

### GetArticlesRequest

Permintaan untuk mengambil artikel pengguna.

| Field     | Type   | Description        |
| --------- | ------ | ------------------ |
| `user_id` | string | ID pengguna target |

---

### GetArticlesResponse

Respons yang berisi daftar artikel milik pengguna.

| Field      | Type                           | Description                   |
| ---------- | ------------------------------ | ----------------------------- |
| `articles` | repeated [`Article`](#article) | Daftar artikel milik pengguna |

---

### Article

Representasi data sebuah artikel.

| Field     | Type   | Description                                |
| --------- | ------ | ------------------------------------------ |
| `id`      | string | ID unik artikel                            |
| `status`  | string | Status artikel (`draft`, `published`, dll) |
| `content` | string | Konten HTML dari artikel                   |
| `user_id` | string | ID pengguna pemilik artikel                |
