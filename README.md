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

- **Backend kept to the standard library.** `net/http`'s Go 1.22+ method-aware routing
  (`GET /health`, `POST /api/calculate`) removes the need for a router dependency for an API this
  small.
- **Operation registry over a `switch`.** `calculator.Evaluate` dispatches through a
  `map[Operation]fn` instead of a growing `switch` statement, so adding a new operation is a
  one-line addition to the map rather than a change to shared dispatch logic. Each operation is a
  small, independently testable function.
- **Business logic isolated from the HTTP layer.** `calculator/` has no knowledge of HTTP/JSON;
  `handlers/` only decodes/validates/encodes and delegates the math. This keeps both layers easy to
  unit test in isolation (see the near-100% coverage on both).
- **Errors as data, not panics.** Invalid input (division by zero, negative square root, zero to a
  negative power, unsupported operation) returns a typed error that the handler turns into a
  `400` with a JSON `{"error": "..."}` body — never a 500 or a crash.
- **Frontend has exactly one `fetch()` call site.** `src/api/calculatorApi.ts` is the only module
  allowed to talk to the network; components never call `fetch` directly, which keeps the UI
  layer trivially testable by mocking a single module.
- **Real calculator UX over a generic form.** The UI mirrors a physical calculator (display +
  circular keypad + chained operations) rather than a plain "enter two numbers" form, per the
  "intuitive UI" requirement. Errors (e.g. divide by zero) render inline in the display itself,
  the same way a real calculator shows `Error`, instead of a separate banner.
- **Square root is unary, add/subtract/multiply/divide/exponent/percentage are binary.** The API
  and UI both model this explicitly: `square_root` ignores `b` and applies immediately; the rest
  participate in the normal "operand → operator → operand → `=`" flow, including chaining.
