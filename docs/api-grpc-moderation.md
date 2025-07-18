# 📄 ModerationService - gRPC API Documentation

## Overview

`ModerationService` adalah gRPC service yang bertugas memoderasi artikel sebelum dipublikasikan. Saat ini hanya memiliki satu method utama:

* `PublishArticle`: Memproses publikasi artikel melalui proses moderasi.

---

## 🛠️ Service: `ModerationService`

### `PublishArticle`

Memvalidasi dan memoderasi artikel sebelum dipublikasikan.

* **Request**: [`PublishArticleRequest`](#publisharticlerequest)
* **Response**: [`PublishArticleResponse`](#publisharticleresponse)

---

## 📦 Messages - ModerationService

### PublishArticleRequest

| Field        | Type   | Description                     |
| ------------ | ------ | ------------------------------- |
| `article_id` | string | ID artikel yang akan dimoderasi |

### PublishArticleResponse

| Field        | Type   | Description                                                    |
| ------------ | ------ | -------------------------------------------------------------- |
| `article_id` | string | ID artikel yang diproses                                       |
| `status`     | string | Status hasil moderasi (`approved`, `rejected`, `pending`, dll) |
