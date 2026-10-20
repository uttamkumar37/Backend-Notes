# 2. HTTP APIs

[Go](01-go-basics.md) → **HTTP** → [Data](03-databases-cache.md)

Goal: follow a request from the network to a validated change and response.

Use the overview to run CRUD. Then study [the network path](#from-url-to-handler), [boundaries](#handler-service-and-repository), [clients](#calling-another-http-service), [security](#authentication-and-authorization), and [API design](#versioning-and-openapi).

## Follow one request

Run [the project](project/README.md), then read:

```text
main.go → api.go → store.go → api_test.go
startup    HTTP      rules      expected behavior
```

Start with standard HTTP to see the request flow. Gin can replace routing later while the store stays testable.

## Define the contract

| Request | Success | Important failure |
| --- | --- | --- |
| `POST /api/v1/investors` | 201 + body + `Location` | 400 invalid input; 415 wrong content type |
| `GET /api/v1/investors/{id}` | 200 + body | 400 malformed ID; 404 missing investor |
| `GET /api/v1/investors?after=0&limit=20` | 200 + ordered array | 400 invalid pagination |
| `PUT /api/v1/investors/{id}` | 200 + replacement | 400 incomplete fields; 404 missing investor |
| `DELETE /api/v1/investors/{id}` | 204, empty body | 404 already missing |

The project trims names, accepts 1–100 Unicode code points, and uppercases a two-letter country label. It does not verify residency. ID belongs to the server.

Use 409 for a business conflict, 500 for an unexpected internal failure, and 503 for temporary unavailability. The current store has no uniqueness conflict. Status codes are part of the contract; an error message with status 200 hides failures. [HTTP semantics](https://www.rfc-editor.org/rfc/rfc9110.html).

## Validate before writing

The implemented decoder requires `application/json`, limits the body to 4096 bytes, and rejects unknown fields, arrays, `null`, malformed JSON, and extra JSON values. Domain validation runs before mutation. Test that rejected requests leave state unchanged.

Standard JSON decoding still accepts duplicate keys, matches field names without case sensitivity, and can replace invalid UTF-8. Rejecting those needs an additional parser/policy. [JSON behavior](https://pkg.go.dev/encoding/json).

Handlers return JSON errors; the standard router supplies plain-text 404/405 errors. Unifying them is a later contract exercise.

## Pagination and safe retries

`after` is the last seen ID, not a page number. Default limit is 20; maximum is 100. Later writes/deletes can change subsequent pages: this is not a cross-request snapshot.

GET should not mutate business state. Repeating the same PUT should leave the same replacement. DELETE can return 204 then 404 and still be idempotent because the resulting state remains deleted.

POST can create duplicates. A client timeout does not prove the write failed. Before retrying writes, design durable idempotency keys that handle concurrent duplicates and changed payloads.

PATCH must distinguish an absent field from an explicit empty/zero value. It is an exercise, not an implemented route.

## Add Gin when the flow is clear

| Standard HTTP | Gin | Your responsibility |
| --- | --- | --- |
| `http.NewServeMux` | `gin.New()` | Routes and contract |
| `r.PathValue("id")` | `c.Param("id")` | Parse and validate |
| JSON decoder | `c.ShouldBindJSON(&input)` | Body bound, business rules, safe errors |
| `r.Context()` | `c.Request.Context()` | Cancellation and I/O deadlines |
| Handler wrapper | Middleware | Auth, recovery, logging, request IDs |

`gin.Default()` adds logging/recovery middleware; `gin.New()` starts without them. `ShouldBindJSON` returns an error for you to handle. Binding does not establish your entire API policy. [Gin binding](https://gin-gonic.com/en/docs/binding/binding-and-validation/).

Authentication establishes identity; authorization decides whether that identity can access this specific record. Add both before serving personal records to users. The local project supplies neither.

## From URL to handler

A URL carries scheme, host, optional port, path, and query. Fragments after `#` are interpreted by the client; they are not part of the HTTP request sent to the server.

For a remote HTTPS call, the client may resolve DNS, establish a connection, negotiate TLS, send the request, and read the response. Reused connections can skip some setup. A failure before a response exists differs from an HTTP error response.

Diagnose the layer: name resolution → connect → TLS → request/response headers → body → business result. “Connection refused” usually means the destination was reached but no suitable listener accepted the connection; a timeout has several possible causes.

## Requests, responses, and headers

| Piece | Example | Purpose |
| --- | --- | --- |
| Method | `POST` | Requested operation |
| Path | `/api/v1/investors/12` | Resource identity |
| Query | `limit=20` | Filtering/pagination options |
| Request header | `Content-Type: application/json` | Format of the supplied body |
| Request header | `Accept: application/json` | Response formats the client accepts |
| Response status | `201` | Outcome category |
| Response header | `Location: /api/v1/investors/12` | Created resource location |
| Body | JSON document | Representation or error detail |

`Content-Type` and `Accept` answer different questions. Header names are case-insensitive. Do not trust arbitrary forwarding headers or request IDs without limits and a trusted-proxy policy.

## A handler without a real port

Complete program — save as a temporary `main.go`, then run `go run main.go`:

```go
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Hello, %s", r.PathValue("name"))
	})
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/hello/Asha", nil)
	mux.ServeHTTP(response, request)
	fmt.Println(response.Code)
	fmt.Println(response.Body.String())
}
```

Output: `200`, then `Hello, Asha`. `ResponseWriter` writes the response; `*http.Request` holds the request. `HandleFunc` adapts a function into HTTP handling. The recorder captures the response in memory.

Set headers before writing the status/body. A first body write implicitly sends status 200 if no final status was written. Once a response is committed, a later error cannot replace it with a new status; design streaming/error handling accordingly.

## Routing and method semantics

Modern `ServeMux` supports method and wildcard patterns. GET patterns also match HEAD; HEAD responses omit the body over the real HTTP transport. A path mismatch differs from a method mismatch, which can return 405 and `Allow`.

Trailing slashes and subtree patterns matter. Treat URL canonicalization, escaped paths, and middleware prefixes as part of routing behavior; test the public URLs rather than guessing from a handler name.

POST is not automatically “create”; it can represent other processing. PUT supplies a replacement at a known resource URL. PATCH supplies a partial change with a documented patch representation. Define these behaviors consistently for your API. [HTTP package](https://pkg.go.dev/net/http).

## Handler, service, and repository

| Layer | Owns | Should avoid |
| --- | --- | --- |
| Handler | HTTP parsing, status, response format | SQL/business decisions scattered through routes |
| Service | Use-case rules, authorization context, transaction boundary | Dependence on a framework context |
| Repository | Storage operations and storage error translation | Writing HTTP responses |

The tiny project combines business rules with its memory store to stay readable. Introduce a service when a use case coordinates several operations. Add interfaces where behavior must be replaced or tested; do not generate an interface for every struct automatically.

**Java bridge:** constructor calls replace bean wiring; a service function should still make its dependencies and failures visible. Avoid porting annotation-heavy designs literally.

## Request and response DTOs

A request DTO describes allowed client input. A response DTO describes the public representation. Domain/storage structs may contain fields that must never be writable or visible to clients.

Never bind a request directly onto a stored object with fields such as tenant ID, role, owner, balance, or internal state. Decide those fields from trusted identity and business rules.

For absent/null/zero values, write examples in the contract. `omitempty` is an encoding rule, not input validation. Use a presence-aware design if PATCH needs to distinguish all three states.

## JSON numbers and body limits

Decoding JSON into `any` can produce floating-point numbers with the compatibility API. Large integer IDs can lose precision when converted through floating point or a JavaScript client. Prefer a typed DTO or an explicit numeric/string contract.

Body-size limits belong before decoding. Also bound decompression, multipart processing, nesting/work where relevant, and downstream response sizes. A small request can still trigger expensive work; body size alone does not establish a compute budget.

The project's decoder demonstrates a pointer target for rejecting `null` and a second decode for rejecting additional JSON values. Reuse the idea with tests rather than adding checks by assumption.

## Middleware and execution order

A middleware wraps the next handler. Think of request order as outer → inner and response/unwind order as inner → outer.

Common concerns: recovery, authentication, authorization context, request IDs, tracing, metrics, and access logs. Ordering changes which failures are observed and which identity is available to later code.

When wrapping `ResponseWriter`, preserve the extra capabilities needed by your handlers, such as flushing or connection upgrading. A wrapper that captures status incorrectly can break streaming or log the wrong status.

## Calling another HTTP service

Complete program using a local test server:

```go
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

func fetch(ctx context.Context, client *http.Client, target string) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 1025))
	if err != nil {
		return "", err
	}
	if len(body) > 1024 {
		return "", fmt.Errorf("response exceeds 1024 bytes")
	}
	return string(body), nil
}

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Mira")
	}))
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	value, err := fetch(ctx, &http.Client{Timeout: 2 * time.Second}, server.URL)
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	fmt.Println(value)
}
```

Output: `Mira`. The extra byte distinguishes “exactly at the limit” from “truncated by our reader.” A non-2xx response does not by itself cause `client.Do` to return an error. Close bodies; consume bounded bodies as appropriate when reusing connections.

Reuse clients/transports across requests and configure them before concurrent use. The example constructs one client for one call; a service would inject its long-lived client. Avoid accepting arbitrary outbound URLs without a deliberate destination policy.

## Timeouts and cancellation

| Control | What it bounds |
| --- | --- |
| Request context | Cooperative request/dependency lifetime |
| Client timeout | Total client request lifetime, including reading its body |
| Connect/TLS/header timeouts | Specific connection/response phases |
| Server header/read timeouts | Incoming network reads |
| Server write timeout | Response network writes under server/protocol semantics |
| Application concurrency limit | How many operations may run at once |

A write timeout does not stop arbitrary CPU work. A context deadline does not cancel an API that never observes context. Avoid treating one timeout setting as a complete overload design.

## Authentication and authorization

Verify identity using the system's chosen session/token mechanism. For signed tokens, validate the expected algorithm, issuer, audience, signature, and time-related claims; receiving a decodable token is not verification. [JWT security practices](https://www.rfc-editor.org/rfc/rfc8725.html).

Authorize the actual operation and object. For multi-tenant data, trusted tenant identity must constrain the repository query, not only a response filter after loading data. A user who knows an investor ID must not gain access merely by changing the path.

Use 401 when authentication is required/invalid according to the contract, and 403 when an identified caller is not allowed. Some APIs deliberately conceal resource existence with 404; make that a consistent policy.

## Cookies, CORS, CSRF, and TLS

Cookies can be attached automatically by browsers; this shapes CSRF protection. Cookie attributes such as `Secure`, `HttpOnly`, and an appropriate `SameSite` policy serve different purposes.

CORS controls which browser origins may read certain responses. It does not prevent all requests and does not authorize a caller. TLS protects transport; it does not validate business input or prove that a user may access a record. [OWASP CSRF guidance](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html).

Behind a reverse proxy, explicitly configure trusted forwarding behavior, original scheme handling, and client IP extraction. Do not let a public client impersonate an internal caller through a header.

## Rate limiting and overload

A rate limit controls requests over time; a concurrency limit controls simultaneous work. Use both where the bottleneck requires them. Define fair keys, burst allowance, retry guidance, and behavior when the limiter's backing store is unavailable.

Per-process limits change when replicas change. A distributed limiter adds its own dependency/failure behavior. Return 429 for a rate-limit policy violation and document any retry information.

Bound queues. An unlimited queue converts overload into memory growth and long latency. Rejecting excess work early can preserve useful service for work already accepted.

## HTTP caching and conditional requests

`Cache-Control` tells clients/intermediaries how a response may be stored and reused. Personal or sensitive responses need an intentional cache policy. `no-cache` permits storage with revalidation; `no-store` directs caches not to store the response.

ETags can identify representation versions. Conditional GET can avoid retransmitting an unchanged representation; conditional writes can prevent overwriting a changed resource when the API defines the precondition.

HTTP response caching differs from an internal database-object cache. Review authorization, variants, and invalidation separately.

## Versioning and OpenAPI

OpenAPI describes routes, parameters, schemas, responses, and authentication requirements. It makes the public contract reviewable; it does not prove the implementation follows it. Test representative contract cases.

Document required versus optional fields, nullability, limits, sorting, error shapes, idempotency, and compatibility rules. Adding a response field is not automatically harmless if strict clients reject unknown fields.

The project uses a simple `{ "error": "..." }` body. A richer API may use standardized problem details, with stable machine-readable type/code and safe human detail. Do not send stack traces or credentials. [Problem details standard](https://www.rfc-editor.org/rfc/rfc9457.html).

## Streaming, uploads, and background work

Streaming commits headers early and requires a plan for mid-stream failure, cancellation, slow readers, and flushing. WebSockets or other upgraded connections need explicit lifecycle management beyond ordinary HTTP shutdown.

For uploads, bound total bytes and parser memory, clean temporary files, and validate content according to use. A filename or MIME header supplied by a client is not a trusted storage path or proof of file type.

For long jobs, consider 202 plus a job resource only after defining how the job is persisted, retried, authorized, and observed. Returning 202 before durable acceptance can silently lose work on restart.

## API debugging table

| Symptom | Check first |
| --- | --- |
| Connection refused | Target address, port, listener, network boundary |
| 404 | Prefix/method/path, then resource existence |
| 405 | Method and `Allow` response header |
| 400 | Parse/validation details and body actually sent |
| 401/403 | Identity validity versus object permission |
| 413 | Configured body limit and proxy limits |
| 415 | Request `Content-Type` |
| 429 | Limiter key, policy, and retry behavior |
| 500 | Correlated internal error; preserve safe client response |
| Slow response | Queue, pool acquisition, query, cache, and downstream timing |

## Practice

1. Complete CRUD with curl. Send an unknown field, missing content type, and `null`; explain the statuses.
2. Add PATCH with tests for absent name, empty name, and missing investor.
3. Replace only HTTP routing with Gin. Reuse the same request/response test cases.

Debug in order: method/path → headers/body → status → validation → store. A wrong route and a missing record can both produce 404 for different reasons.

Check yourself: why 201? What makes DELETE idempotent? Why can retrying POST after a timeout duplicate data?

Additional challenges: unify router/handler error shapes; add ETag-based conditional update; add a body-size regression test; prove that a second fictional tenant cannot read the first tenant's investor. Keep each change small enough to review independently.

## Quick answers

| Question | Answer |
| --- | --- |
| Does `client.Do` fail on HTTP 500? | It can return a valid response with nil error; inspect the status. |
| Why close a response body? | Release owned resources and support connection reuse appropriately. |
| Does `omitempty` validate required input? | No; it controls encoding behavior. |
| Can CORS replace authorization? | No; it is a browser response-access policy. |
| Does a timed-out POST definitely fail? | No; the server may already have committed the effect. |
| Does a body limit stop expensive queries? | No; compute/concurrency/query budgets are separate controls. |
| Can a handler rewrite status after streaming starts? | A committed status cannot be replaced by a new final status. |
| Should Gin context become the service API? | Prefer domain inputs plus standard context across framework boundaries. |
