# File Upload API Documentation

## Base URL

```
http://<server-address>/
```

---

## Endpoints

### 1. Upload a File

**URL:** `/`
**Method:** `POST`
**Authentication:** Required (`Authorization` header)

**Headers:**

```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Body Parameters:**

* `file` — The file to upload
* `name` — The username of the uploading user (parsed from the form)

**Response:**

```json
{
  "status": "success",
  "message": "Upload complete, download your file here: http://<host>/<short_code>",
  "url": "http://<host>/<short_code>"
}
```

**Errors:**

* `401 Unauthorized` — Invalid or missing token
* `404 User not found` — User ID does not exist
* `400 Bad Request` — Multipart parsing error
* `500 Internal Server Error` — File save or database error

---

### 2. Get All User Uploads

**URL:** `/uploads/`
**Method:** `GET`
**Authentication:** Required (`Authorization` header)

**Headers:**

```
Authorization: Bearer <token>
```

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "short_code": "abcd",
      "filename": "example.txt"
    }
  ]
}
```

**Errors:**

* `401 Unauthorized` — Invalid or missing token
* `500 Internal Server Error` — Database error

**Notes:**

* `POST` requests to this endpoint are not allowed (`405 Method Not Allowed`).

---

### 3. Download a File via Short Code

**URL:** `/<short_code>`
**Method:** `GET`
**Authentication:** Not required

**Response:**

* Serves the file directly if the short code exists.

**Errors:**

* `404 Not Found` — Invalid or missing short code

**Notes:**

* Short codes must be exactly 4 characters.

---

### 4. Serve Web Application

**URL:** `/app/<path>`
**Method:** `GET`
**Authentication:** Not required

**Response:**

* Serves static files from the directory specified in the `WEBAPP_PATH` environment variable.

---

## Data Models

### User

```go
type User struct {
    ID           int
    Name         string
    PasswordHash string
    SessionToken string
}
```

### Upload

```go
type Upload struct {
    ID        int
    UserID    int
    FileName  string
    FilePath  string
    ShortCode string
}
```

---

### Authentication

* All POST requests (except `/app/`) and `/uploads/` GET requests require a Bearer token in the `Authorization` header.

