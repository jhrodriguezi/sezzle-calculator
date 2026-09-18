# Sezzle Calculator API

A minimal Go HTTP API that performs basic arithmetic calculations.

## Running

```bash
go run .
```

The server listens on `:8080` by default. Set the `PORT` environment variable to override.

## Endpoints

### `GET /health`

Health check endpoint used for readiness/liveness checks.

**Request:** No body required.

**Response:** `200 OK`

```json
{
  "status": "ok"
}
```

### `POST /api/calculate`

Performs an arithmetic operation on two numbers.

**Request body:**

| Field       | Type   | Required | Description                                              |
|-------------|--------|----------|------------------------------------------------------------|
| `operation` | string | yes      | One of `add`, `subtract`, `multiply`, `divide`             |
| `a`         | number | yes      | First operand                                              |
| `b`         | number | yes      | Second operand                                              |

Example:

```json
{
  "operation": "add",
  "a": 2,
  "b": 3
}
```

**Success response:** `200 OK`

```json
{
  "result": 5
}
```

**Error responses:** `400 Bad Request`

Returned when the request body is malformed, the `operation` is unsupported, or a division by zero is attempted.

```json
{
  "error": "unsupported operation: modulo"
}
```

```json
{
  "error": "division by zero"
}
```

```json
{
  "error": "invalid request body"
}
```
