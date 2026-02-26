# URL Shortener API Specification

## Overview

Base URL: `https://api.example.com/v1`

The URL Shortener API allows users to create shortened URLs, track click statistics, and manage shortened URLs.

---

## Endpoints

### 1. Redirect to Original URL

**Endpoint:** `GET /:shortCode`

**Description:** Redirects to the original URL associated with the provided short code. Registers a click asynchronously.

**Path Parameters:**

| Parameter   | Type   | Description                              | Constraints        |
|-------------|--------|------------------------------------------|--------------------|
| shortCode   | string | The shortened URL code                   | Must be 7 characters |

**Response (301 - Moved Permanently):**

The server responds with an HTTP redirect to the original URL. The `Location` header contains the original URL.

**Error Responses:**

| Status Code | Description                           |
|-------------|---------------------------------------|
| 400         | Invalid shortCode format              |
| 404         | URL not found                         |
| 410         | URL expired or has been deleted       |
| 500         | Internal server error                 |

**Error Response Body:**
```json
{
  "error": "error message"
}
```

---

### 2. Create Short URL

**Endpoint:** `POST /shorten`

**Description:** Creates a new shortened URL.

**Request Body:**

| Field          | Type    | Required | Description                          |
|----------------|---------|----------|--------------------------------------|
| url            | string  | Yes      | The original URL to shorten          |
| user_id        | string  | No       | User identifier for tracking         |
| expires_in_days| uint32 | No       | Days until the URL expires            |

**Example Request:**
```json
{
  "url": "https://www.example.com/some-long-url",
  "user_id": "user123",
  "expires_in_days": 30
}
```

**Success Response (201 - Created):**

```json
{
  "short_code": "abc1234",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

| Field     | Type     | Description                          |
|-----------|----------|--------------------------------------|
| short_code| string  | The generated short code (7 chars)  |
| expires_at| string  | Expiration timestamp (ISO 8601)     |

**Error Responses:**

| Status Code | Description                           |
|-------------|---------------------------------------|
| 400         | Invalid request body or missing URL  |
| 500         | Internal server error                 |

**Error Response Body:**
```json
{
  "error": "error message"
}
```

---

### 3. Get URL Statistics

**Endpoint:** `GET /stats/:shortCode`

**Description:** Retrieves overall statistics for a shortened URL.

**Path Parameters:**

| Parameter   | Type   | Description                              | Constraints        |
|-------------|--------|------------------------------------------|--------------------|
| shortCode   | string | The shortened URL code                   | Must be 7 characters |

**Success Response (200 - OK):**

```json
{
  "short_code": "abc1234",
  "total_clicks": 1500,
  "unique_visitors": 450,
  "referers": [
    "https://google.com",
    "https://twitter.com"
  ],
  "last_clicked_at": "2024-06-15T14:30:00Z"
}
```

| Field           | Type    | Description                          |
|-----------------|---------|--------------------------------------|
| short_code      | string  | The short code                       |
| total_clicks    | integer | Total number of clicks               |
| unique_visitors | integer | Number of unique visitors            |
| referers        | array   | List of referring URLs (may be empty)|
| last_clicked_at | string  | Timestamp of last click (ISO 8601)  |

**Error Responses:**

| Status Code | Description                           |
|-------------|---------------------------------------|
| 400         | Invalid shortCode format              |
| 500         | Internal server error                 |

**Error Response Body:**
```json
{
  "error": "error message"
}
```

---

### 4. Get Daily Statistics

**Endpoint:** `GET /stats/:shortCode/:date`

**Description:** Retrieves statistics for a specific date.

**Path Parameters:**

| Parameter   | Type   | Description                              | Constraints        |
|-------------|--------|------------------------------------------|--------------------|
| shortCode   | string | The shortened URL code                   | Must be 7 characters |
| date        | string | The date in YYYY-MM-DD format             | Valid date string  |

**Example:** `GET /stats/abc1234/2024-06-15`

**Success Response (200 - OK):**

```json
{
  "short_code": "abc1234",
  "date": "2024-06-15",
  "total_clicks": 75,
  "unique_visitors": 30,
  "referers": [
    "https://google.com"
  ]
}
```

| Field           | Type    | Description                          |
|-----------------|---------|--------------------------------------|
| short_code      | string  | The short code                       |
| date            | string  | The date (YYYY-MM-DD)                |
| total_clicks    | integer | Number of clicks on this date       |
| unique_visitors | integer | Number of unique visitors           |
| referers        | array   | List of referring URLs (may be empty)|

**Error Responses:**

| Status Code | Description                           |
|-------------|---------------------------------------|
| 400         | Invalid shortCode or date format      |
| 500         | Internal server error                 |

**Error Response Body:**
```json
{
  "error": "error message"
}
```

---

### 5. Delete Short URL

**Endpoint:** `DELETE /:shortCode`

**Description:** Soft deletes a shortened URL by marking it as inactive.

**Path Parameters:**

| Parameter   | Type   | Description                              | Constraints        |
|-------------|--------|------------------------------------------|--------------------|
| shortCode   | string | The shortened URL code                   | Must be 7 characters |

**Success Response (200 - OK):**

Empty response body. The server returns status 200 on successful deletion.

**Error Responses:**

| Status Code | Description                           |
|-------------|---------------------------------------|
| 400         | Invalid shortCode format              |
| 500         | Internal server error                 |

**Error Response Body:**
```json
{
  "error": "error message"
}
```

---

## Common Data Types

### Error Response

All error responses follow this format:

```json
{
  "error": "error message"
}
```

### Date Format

All dates and timestamps use ISO 8601 format:
- Date: `YYYY-MM-DD`
- Timestamp: `YYYY-MM-DDTHH:MM:SSZ` (UTC)

### Short Code Format

- Length: Exactly 7 characters
- Format: Alphanumeric string

---

## HTTP Status Codes Summary

| Status Code | Meaning                                    |
|-------------|-------------------------------------------|
| 200         | OK - Request succeeded                    |
| 201         | Created - Resource successfully created   |
| 301         | Moved Permanently - Redirect              |
| 400         | Bad Request - Invalid input               |
| 404         | Not Found - Resource not found            |
| 410         | Gone - Resource expired or deleted        |
| 500         | Internal Server Error                     |

---

## Notes for Frontend Developers

1. **Redirect Handling:** When accessing a short URL, the browser will automatically redirect. Ensure your application handles the 301 redirect response properly.

2. **Optional Fields:** Some response fields may be null or omitted (e.g., `expires_at`, `referers`). Always handle optional fields gracefully.

3. **Validation:** The short code must be exactly 7 characters. Dates must be in `YYYY-MM-DD` format.

4. **Async Analytics:** Click tracking happens asynchronously, so real-time stats may have a slight delay.
