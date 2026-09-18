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

Performs an arithmetic operation on one or two numbers.

**Request body:**

| Field       | Type   | Required | Description                                                                                     |
|-------------|--------|----------|---------------------------------------------------------------------------------------------------|
| `operation` | string | yes      | One of `add`, `subtract`, `multiply`, `divide`, `exponentiate`, `square_root`, `percentage`        |
| `a`         | number | yes      | First operand                                                                                      |
| `b`         | number | no       | Second operand. Ignored for `square_root` (unary operation)                                       |

Operation semantics:

| Operation      | Result                          |
|----------------|----------------------------------|
| `add`          | `a + b`                          |
| `subtract`     | `a - b`                          |
| `multiply`     | `a * b`                          |
| `divide`       | `a / b` (error if `b` is 0)       |
| `exponentiate` | `a ^ b`                           |
| `square_root`  | `√a` (error if `a` is negative)   |
| `percentage`   | `a` percent of `b`, i.e. `(a / 100) * b` |

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

Returned when the request body is malformed, the `operation` is unsupported, a division by zero is attempted, a negative number's square root is requested, or the operation otherwise produces an invalid result (e.g. `0` raised to a negative power).

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
  "error": "cannot take the square root of a negative number"
}
```

```json
{
  "error": "invalid request body"
}
```
