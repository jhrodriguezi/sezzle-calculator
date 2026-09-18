# Sezzle Calculator

A full-stack calculator: a Go REST API for arithmetic (basic + advanced operations) and a
React/TypeScript UI that looks and behaves like a real calculator.

```
.
├── backend/   Go HTTP API (see backend/README.md for full API reference)
└── frontend/  React + Vite + TypeScript UI (see frontend/README.md for details)
```

## Setup

### Prerequisites

- Go 1.25+
- Node.js 22+
- (Optional) Docker / Docker Compose

### Backend

```bash
cd backend
go run .
```

Listens on `:8080` by default (override with the `PORT` env var).

### Frontend

```bash
cd frontend
cp .env.example .env   # sets VITE_API_BASE_URL=http://localhost:8080
npm install
npm run dev
```

Opens on `http://localhost:5173` and talks to the backend at `VITE_API_BASE_URL`.

### Run both with Docker Compose

```bash
docker compose up --build
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000`

## API examples

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 2, "b": 3}'
# {"result":5}

curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "square_root", "a": 81}'
# {"result":9}

curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 1, "b": 0}'
# 400 Bad Request: {"error":"division by zero"}
```

Full endpoint reference (all supported operations, request/response shapes, error cases) is in
[`backend/README.md`](backend/README.md).

## Tests & coverage

### Backend

```bash
cd backend
go test ./...                              # run tests
go test ./... -coverprofile coverage.out   # with coverage profile
go tool cover "-func=coverage.out"         # per-function summary
go tool cover "-html=coverage.out"         # HTML report
```

Current coverage: `calculator` package ~95%, `handlers` package 100% (the `main` package is just
server wiring and is intentionally untested).

### Frontend

```bash
cd frontend
npm run test            # run tests (Vitest + React Testing Library)
npm run test:coverage   # with coverage report (text + HTML in coverage/)
```

## Design decisions

- **Scope.** "Basic and advanced arithmetic operations" mapped to the 4 required ops
  (add/subtract/multiply/divide) plus the 3 listed-as-optional ones (exponentiation, square root,
  percentage). An expression parser (tokenizer → Shunting-yard → evaluator) was considered and
  deliberately dropped — it wasn't asked for, and adds real complexity (precedence, parentheses,
  unary-minus) that trades a complete, well-tested core for a half-finished extra. Chaining
  (`4 + 7 =`, then `* 2 =`) still works, but lives entirely in frontend UI state — each `=` is one
  independent `a/b` request to the backend.
- **Backend kept to the standard library.** `net/http`'s Go 1.22+ method-aware routing
  (`GET /health`, `POST /api/calculate`) removes the need for a router dependency, ORM, or HTTP
  client for an API this small.
- **One endpoint, fixed `a`/`b` fields.** A single `POST /api/calculate` takes `operation` + `a`/`b`
  rather than one endpoint per operation or a variable-length operand list — one validation/error
  path, one frontend call shape, and a new operation is one registry entry (below), not a new
  route. The unary operation (`square_root`) simply ignores `b` (JSON decodes a missing field to
  `0`); splitting unary/binary into separate endpoints was considered and set aside as unnecessary
  overhead for a single unary op out of seven.
- **Operation registry over a growing `switch`.** `calculator.Evaluate` dispatches through a
  `map[Operation]fn`, so adding an operation is a one-line map entry, not a change to shared
  dispatch logic. Each operation is a small, independently testable function.
- **Business logic isolated from the HTTP layer.** `calculator/` has no knowledge of HTTP/JSON;
  `handlers/` only decodes/validates/encodes and delegates the math — no service/repository layer,
  since there's no persistence to justify one. Keeps both layers easy to unit test in isolation
  (near-100% coverage on both).
- **Errors as data, not panics.** Invalid input (division by zero, negative square root, zero to a
  negative power, unsupported operation) returns a typed error that the handler turns into a `400`
  with a JSON `{"error": "..."}` body — never a 500, never a `recover()`.
- **Frontend has exactly one `fetch()` call site.** `src/api/calculatorApi.ts` is the only module
  allowed to talk to the network; components never call `fetch` directly, keeping the UI layer
  trivially testable by mocking a single module. No global state library or router — a single
  screen with `useState` doesn't need one.
- **Real calculator UX over a generic form.** The UI mirrors a physical calculator (display +
  keypad + chained operations) rather than a plain "enter two numbers" form. Errors render inline
  in the display, the same way a real calculator shows `Error`, instead of a separate banner.
- **Assumptions.** Percentage is `(a / 100) * b` ("a percent of b") since the shape doesn't imply a
  semantic; negative square-root input returns a `400` rather than a complex result.
