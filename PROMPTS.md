# AI Prompts Used

This log combines two sessions: architecture/planning discussion (before and alongside implementation) and the
hands-on coding session. Purely mechanical prompts (commit message wording, branch housekeeping)
are left out.

## Architecture & design discussion

1. What is the best approach then? *(followed by the full assignment text
   — used to work through scope, API shape, and structure before writing
   code)*
2. If we want a square root operation, we'd have to pass the value into
   one of the two operand fields — does that make sense with the fixed
   `a`/`b` shape, or is there a cleaner way to handle a unary operation
   here?
3. From a api-design perspective, does it make sense to force the
   consumer to send two operands? Would it be better to separate unary
   operations into its own endpoint?
4. Lets adjust that part because at end I went with this approach
   *(providing the actual final Go structs and asking for the design
   rationale to be reconciled with what was actually built, not what was
   originally proposed)*
5. How would you evolve that to include the advanced arithmetic
   operations?

## Implementation session

6. Keep track of every prompt I send in this session — I'll want the full
   list summarized at the end.
7. Scaffold a Go API in a `backend/` folder using only the standard
   library (`net/http`) — no third-party router. Go 1.22's method-based
   routing is enough for this scope.
8. Drop the `cmd/`, `internal/`, `service/`, `repository/` layout —
   restructure it into a flat `calculator/` (pure logic) and `handlers/`
   (HTTP) split instead.
9. In the README, document the single `POST /api/calculate` endpoint:
   the request shape (`operation`, `a`, `b`), the response shape
   (`result`), and the error shape (`error`), with one example per
   operation.
10. Build the frontend in its own folder, calling that endpoint through
    one dedicated API module — no other component should call `fetch`
    directly.
11. Style it like an actual calculator: a display plus a keypad layout,
    not a plain form with two number inputs.
12. Show errors inline inside the display itself, the way a physical
    calculator shows `Error`, instead of a separate banner or toast.
13. Add the advanced operations — exponentiation, square root, and
    percentage — into the same operation dispatch, keeping square root
    unary and reading only `a`.
14. Redesign how the sqrt, exponent, and percentage controls look in the
    keypad — right now they don't fit visually with the basic operator
    buttons.