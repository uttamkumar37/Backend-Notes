# 3. Databases and caching

[HTTP](02-http-api.md) → **Data** → [Concurrency](04-concurrency.md)

Goal: keep correct data across restarts, concurrent requests, and dependency failures.

The supplied project uses memory. Everything below is a sequence for extending it.

Use [schema design](#schema-and-constraints), [Go SQL examples](#a-context-aware-query), [concurrent writes](#optimistic-and-pessimistic-concurrency), [query plans](#reading-a-query-plan), and [cache behavior](#cache-keys-and-serialization) as separate reading stops.

## PostgreSQL first

Learn in this order: SQL/schema → migrations → parameterized CRUD → transactions → indexes → pool behavior.

Use `database/sql` with a driver, or native `pgx` with `pgxpool`; understand which API your service uses. A pool reuses a finite number of connections. Create it once at startup, verify connectivity, pass it to the repository, and close it after active work finishes. [Go database guide](https://go.dev/doc/database), [pgx docs](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool).

Introduce a repository interface when replacing the memory store. Its I/O methods should accept `context.Context`; pass it into database operations. Put the interface with its consumer and keep HTTP types outside the database package.

For investors, define a schema with an identity ID, required name, two-letter country label, and constraints that protect valid rows. Version schema changes in migrations. Test against a fresh database; handle upgrades before declaring the service ready.

**Parameterized query shape** — pass values separately through the driver:

```sql
SELECT id, name, country
FROM investors
WHERE id > $1
ORDER BY id
LIMIT $2;
```

Never concatenate untrusted values into SQL. Parameters represent values, not identifiers; whitelist permitted sort columns. Close result rows, check iteration errors, and map “no rows” to domain not-found behavior.

## Transactions protect business rules

Suppose a fictional transfer debits one balance and credits another. Both changes must commit together or both must roll back. Validate the business invariant inside the transaction, and call all related operations through that transaction handle. [Go transaction guidance](https://go.dev/doc/database/execute-transactions).

PostgreSQL's default Read Committed isolation takes a new snapshot for each statement. It does not make a multi-step read/check/write invariant safe automatically. Use appropriate row locks, atomic conditional updates, constraints, or a suitable isolation level. Serializable transactions can require retrying the **whole transaction**. [PostgreSQL isolation](https://www.postgresql.org/docs/current/transaction-iso.html).

For lost updates, consider a version column: update only when the supplied version matches, then increment it. Zero affected rows can indicate a conflict. Keep network calls outside long-held database transactions; use an outbox when a committed write must reliably lead to a later event.

## Measure queries and pools

- Inspect actual query patterns before adding indexes. Indexes cost storage and write work.
- Use `EXPLAIN` to inspect a plan. `EXPLAIN ANALYZE` executes the statement; use controlled data.
- Pool size is a shared database budget: replicas × per-replica maximum, plus migrations/admin headroom, must fit server capacity.
- Separate pool-acquisition wait from query duration. A deadline can expire before SQL reaches the server.

Start small and compare measured latency, connection use, and query plans. [PostgreSQL indexes](https://www.postgresql.org/docs/current/indexes.html).

## Add Valkey only after correct persistence

Cache-aside read: cache hit → return; miss → read PostgreSQL → cache with TTL → return. On a cache outage, a cache-optional read may fall back to the database **within a bounded load budget**.

Write to PostgreSQL first, then invalidate the relevant cache key. Test the stale-reader race: an earlier read can populate old data after invalidation. TTL bounds a stale window; strict consistency may require versions or a stronger design. A cache is not the authority for balances.

Track hit/miss counts, cache failures, and database load. Limit simultaneous cache fills to avoid a stampede. Cache keys need entity/version and, where applicable, tenant boundaries. A cache failure must not bypass authorization.

Sentinel provides monitoring/failover; Cluster provides partitioning. Asynchronous replication means failover can lose acknowledged writes. A distributed lock needs ownership, expiry handling, and an explicit correctness model; a simple cache lock does not make a financial operation safe. [Valkey Sentinel](https://valkey.io/topics/sentinel/).

## Choose the storage model

Start with the operations and invariants. Which data must survive restart? Which rows change together? Which reads need the latest committed value? Which fields identify a tenant or owner?

Normalize data when separate entities have their own identity and constraints. Denormalize deliberately for a measured read need, with a plan to keep copies consistent. A JSON column is useful for some variable attributes; it does not remove the need to model searchable fields and invariants.

## Schema and constraints

Illustrative PostgreSQL schema for a **future** tenant-aware repository:

```sql
CREATE TABLE investors (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id text NOT NULL,
    name text NOT NULL CHECK (char_length(name) BETWEEN 1 AND 100),
    country text NOT NULL CHECK (country ~ '^[A-Z]{2}$'),
    version bigint NOT NULL DEFAULT 1 CHECK (version > 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX investors_tenant_id_id_idx ON investors (tenant_id, id);
```

`PRIMARY KEY` provides identity. `NOT NULL`, `CHECK`, `UNIQUE`, and foreign keys protect different rules. Application validation produces useful client feedback; database constraints protect writes from every caller.

This schema checks stored length/shape; the application still owns Unicode whitespace normalization and the allowed-country policy. A two-letter label is not proof of real residency. The current in-memory project has no tenant field.

Name important constraints so storage errors can be mapped deliberately. Avoid leaking raw SQL/constraint details directly to the public API. [PostgreSQL constraints](https://www.postgresql.org/docs/current/ddl-constraints.html).

## Null, zero, and absent data

SQL NULL represents a missing/unknown value; it is not an empty string or zero. Comparisons involving NULL follow three-valued logic. Use `IS NULL` rather than `= NULL`.

Choose a Go representation that preserves the contract: pointer, `sql.NullString`/other nullable types, or the driver's supported type. Do not silently convert “unknown amount” to zero if the distinction matters.

`COALESCE` is appropriate only when the fallback means the same thing to the business. Decide how nullable data appears in JSON independently from how the driver scans it.

## database/sql versus native pgx

`*sql.DB` is a concurrency-safe pool handle, not a single connection. Native `pgxpool.Pool` is a different API with its own types/options. `pgx` can also provide a `database/sql` driver through its adapter.

Choose one style for a small service and follow its API consistently. Check startup connectivity explicitly; constructing a handle does not always prove a live authenticated connection. Close it once the application has drained users of it. [database/sql](https://pkg.go.dev/database/sql).

## A context-aware query

Package example — compiles as a Go package; execution needs a configured PostgreSQL-compatible `*sql.DB` and the schema above:

```go
package notes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrInvestorMissing = errors.New("investor not found")

func FindName(ctx context.Context, db *sql.DB, tenant string, id int64) (string, error) {
	const query = `SELECT name FROM investors WHERE tenant_id = $1 AND id = $2`
	var name string
	err := db.QueryRowContext(ctx, query, tenant, id).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvestorMissing
	}
	if err != nil {
		return "", fmt.Errorf("find investor: %w", err)
	}
	return name, nil
}
```

The tenant value must come from trusted authorization context, not an unverified body field. Values are parameters; the SQL structure is constant. `Scan` exposes a query-row error, including no matching row.

For multi-row queries: check the query error → defer `rows.Close()` → iterate `rows.Next()` → scan/check each row → check `rows.Err()`. Leaving rows open can retain a pool connection. [Go query guidance](https://go.dev/doc/database/querying).

## Transaction lifecycle

The owner begins a transaction, performs all related operations through it, commits on success, and rolls back when leaving unsuccessfully. A rollback defer is commonly retained even after a successful commit; cleanup after completion does not undo the committed transaction.

Do not use the pool for one step that belongs in the active transaction. That step may use another connection and commit independently. Do not hold transactions open while waiting for an HTTP call or human input.

Errors can occur during begin, a statement, scanning, or commit. A connection loss around commit may leave the client uncertain whether the server committed; blindly repeating the business operation can duplicate effects.

## A guarded transfer example

Package example for a separate fictional `accounts(id, balance_minor)` table. It illustrates transaction flow; the table, driver, integration tests, and retry/idempotency policy are not supplied:

```go
package notes

import (
	"context"
	"database/sql"
	"fmt"
)

func Transfer(ctx context.Context, db *sql.DB, from, to, amount int64) error {
	if from <= 0 || to <= 0 || from == to || amount <= 0 {
		return fmt.Errorf("invalid transfer")
	}
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin transfer: %w", err)
	}
	defer tx.Rollback()
	debit, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance_minor = balance_minor - $1
		 WHERE id = $2 AND balance_minor >= $1`, amount, from)
	if err != nil {
		return fmt.Errorf("debit: %w", err)
	}
	changed, err := debit.RowsAffected()
	if err != nil {
		return fmt.Errorf("debit result: %w", err)
	}
	if changed != 1 {
		return fmt.Errorf("source missing or insufficient balance")
	}
	credit, err := tx.ExecContext(ctx,
		`UPDATE accounts SET balance_minor = balance_minor + $1 WHERE id = $2`, amount, to)
	if err != nil {
		return fmt.Errorf("credit: %w", err)
	}
	changed, err = credit.RowsAffected()
	if err != nil {
		return fmt.Errorf("credit result: %w", err)
	}
	if changed != 1 {
		return fmt.Errorf("destination missing")
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transfer: %w", err)
	}
	return nil
}
```

The guarded debit prevents a negative result from this statement. The transaction prevents a debit being left committed if the destination update fails. Real use also needs schema constraints, amount bounds/currency rules, stable error categories, authorization, operation identity, and concurrency/retry tests.

## Optimistic and pessimistic concurrency

An optimistic update carries the version that was read:

```sql
UPDATE investors
SET name = $1, version = version + 1
WHERE tenant_id = $2 AND id = $3 AND version = $4
RETURNING id, version;
```

No returned row can mean a stale version or missing/inaccessible investor; map it according to the API's disclosure policy. The client must resolve a conflicting edit rather than silently overwriting it.

Pessimistic coordination may use row locks such as `SELECT ... FOR UPDATE` inside a transaction. Hold locks briefly and acquire multiple resources in a consistent order. Deadlocks still need a recovery policy. [PostgreSQL locking](https://www.postgresql.org/docs/current/explicit-locking.html).

## Isolation and retries

| Isolation | What to reason about |
| --- | --- |
| Read Committed | Successive statements can observe different committed snapshots |
| Repeatable Read | A stable transaction snapshot does not protect every cross-row invariant |
| Serializable | Transactions may be rejected to preserve serializable outcomes |

Retry only classified retryable failures, rerunning the entire transaction with a bounded attempt/time budget. Do not retry an invalid input, a uniqueness violation that represents a business conflict, or an uncertain side effect indiscriminately.

Keep the transaction body replayable. Sending a notification inside it may duplicate the notification during retries, even when the database rolls back.

## Pool sizing and saturation

Example planning budget: a server permits 100 connections; reserve 20 for migrations/admin/other users; four API replicas could each receive at most 20 of the remaining 80. Autoscaling and overlapping rollout replicas change that arithmetic.

Pool size is not a target for maximizing throughput. More simultaneous queries can overload the database and increase latency. Measure active/idle connections, acquisition waits, query duration, and cancellations.

Limit lifetime/idle time according to network/server behavior. Distinguish “pool is full because traffic increased” from “connections never return because rows/transactions leaked.” Raising the maximum can hide the leak temporarily.

## Indexes and query shape

Design an index for a real predicate and ordering. An index on `(tenant_id, id)` fits “one tenant, IDs after cursor, ordered by ID.” Column order matters to useful access patterns.

An index adds write/storage cost. A sequential scan can be appropriate for a small table or a query returning much of it. Avoid declaring every scan without an index a failure.

Fetch only needed columns, avoid unbounded lists, and examine N+1 access patterns. One request that loads a list and then issues one query per item can amplify latency and pool pressure.

## Reading a query plan

Inspect estimated versus actual rows, chosen scan/join, filtering, loops, sorting, buffers, and total execution time. A large estimate error can point to stale statistics or a data-distribution assumption.

`EXPLAIN ANALYZE` executes the statement. A rollback can protect ordinary transactional changes during an experiment, but not every side effect is undone, such as sequence advancement. Use controlled data and a clear plan. [EXPLAIN guide](https://www.postgresql.org/docs/current/using-explain.html).

Measure with representative row counts and parameters. A fast query on ten fictional rows does not establish production capacity.

## Migrations and rollback

Keep ordered migrations in version control. Test a fresh schema and an upgrade from the previous supported schema. Coordinate migration ownership rather than letting every new replica race to change tables.

Use expand → migrate/backfill → switch callers → contract for incompatible changes. An application rollback cannot safely restore an older binary after its required column has been dropped.

Large backfills need batches, progress tracking, resumability, and lock/load awareness. Schema compatibility is a release design problem, not only a migration-tool command.

## Cache keys and serialization

A useful cache key identifies entity, schema/version, and tenant/authorization boundary where applicable. Example: `investor:v1:tenant-demo:42`. Do not let one tenant's key alias another's record.

Choose what is cached: a domain snapshot, public DTO, or expensive computed result. Record serialization version and TTL behavior. Cache corruption/old formats should become observable misses or explicit failures according to policy, not panic the handler.

Cache-aside applies to disposable copies. If the application uses Valkey for durable business state, it needs a different durability/recovery design than an optional cache.

## Valkey data types

| Type | Typical use |
| --- | --- |
| String | Serialized object, counter, small value |
| Hash | Fields grouped under one key |
| Set | Unique membership |
| Sorted set | Ordered scores/ranking/window bookkeeping |
| List | Ordered push/pop data; delivery semantics need care |
| Stream | Append-oriented events with consumer-group features |

Select the type from required operations and consistency guarantees, not only familiarity. Pipelining reduces round trips; it does not automatically make commands an atomic transaction. [Valkey data types](https://valkey.io/topics/data-types/).

## TTL and key commands

Illustrative commands for an already-running disposable lab instance:

```text
SET investor:demo "fictional" EX 60
GET investor:demo
TTL investor:demo
DEL investor:demo
```

TTL reports seconds remaining, `-1` for an existing key without expiry, and `-2` for a missing key. Prefer setting value and expiry together when both are required; separate writes can leave a permanent key if the second command fails. [SET options](https://valkey.io/commands/set/).

Use bounded `SCAN` iteration rather than treating a full keyspace operation as harmless on a busy instance. Exact maintenance commands depend on the environment and workload.

## go-redis client behavior

`go-redis/v9` is a commonly used Go client. Create/reuse a client, configure connection/time budgets, verify connectivity when required, and close it during cleanup. Separate a missing-key result from a transport/authentication error.

`NewClient` targets a configured server; `NewFailoverClient` adds Sentinel discovery configuration. A pool belongs to the client; creating a client per HTTP request defeats reuse and multiplies connection churn. Pin the chosen module version when implementing. [Go client guide](https://redis.io/docs/latest/develop/clients/go/).

## Stampedes, negative caching, and stale data

A hot key expiring can send many identical requests to the database. Options include per-key request coalescing, bounded refresh workers, staggered TTLs, and carefully defined stale-while-revalidate behavior.

Negative caching stores “missing” briefly. It reduces repeated misses but must not conceal a newly created record longer than the contract allows. Do not cache authentication/authorization failures as ordinary global entity misses.

A failed invalidation after a successful database commit leaves a stale cache. Define bounded staleness, repair, metrics, and possibly durable invalidation delivery. A successful invalidation still does not eliminate an old read racing to refill the key.

## Sentinel and distributed locks

Sentinel quorum concerns failure agreement; failover authorization also depends on a majority of Sentinels. Test discovery from the actual client network because an advertised address reachable by containers may not be reachable from the macOS host.

Use a unique owner token for a lease-style lock and release it only if ownership still matches. A paused worker can resume after its lease expires; protecting an external resource may require fencing/version checks. Define the required safety model before choosing a lock algorithm.

## Outbox and event delivery

When a database write must lead to a later notification, store an outbox record in the **same transaction**. A worker publishes it and records progress. A crash between publication and marking completion can cause duplicate delivery.

Consumers therefore need idempotent handling or a deduplication contract. An outbox provides a reliable bridge from committed state to later delivery; it does not create universal exactly-once effects across arbitrary systems.

## Backup, retention, and deletion

Replication helps availability; it can replicate an accidental deletion too. A backup is useful only when restoration is tested against a documented recovery objective.

Soft deletion keeps rows and changes every uniqueness/query/retention decision. Hard deletion must account for dependencies, caches, events, and backups. Choose the model from the actual product/privacy requirements.

## Data debugging table

| Symptom | Separate these causes |
| --- | --- |
| Connection refused | Wrong address/port, stopped listener, network boundary |
| Authentication failed | Wrong identity/password, database/role mismatch |
| Relation missing | Wrong database/schema/search path, unapplied migration |
| Deadline exceeded | Pool wait, lock wait, query execution, network response |
| Increasing pool wait | Higher load, slow statements, leaked rows/transactions |
| Duplicate data | Missing uniqueness/idempotency or replayed write |
| Stale cache | TTL, failed invalidation, old reader refill, wrong key |
| Cache outage overloads DB | Unbounded fallback and missing load protection |
| Failover never reaches client | Discovery configuration, network reachability, authentication |

## Practice

1. Replace memory with PostgreSQL. Prove create → restart API → get still works. Test invalid input and rollback against an isolated real database.
2. Run two concurrent updates. Demonstrate the lost-update problem, then protect the invariant and test the conflict.
3. Add cache-aside. Demonstrate miss → hit → update → fresh read, TTL expiration, and a stopped-cache fallback without overwhelming the database.

Debug by separating connectivity, authentication, migration/schema, query, and pool-wait failures. Read the cause before increasing timeouts.

Check yourself: why is a pool finite? Why isn't a transaction alone enough for every race? How can a stale value reappear after cache invalidation?

Further drills: add optimistic version conflicts; run a migration upgrade test; find and remove an N+1 query; measure pool waits under a bounded load; demonstrate duplicate outbox delivery and idempotent consumption. Record evidence from real isolated dependencies before marking those skills complete.

## Quick answers

| Question | Answer |
| --- | --- |
| Is `sql.DB` one connection? | No; it manages a pool. |
| Does an HTTP request need its own pool? | No; reuse the application-owned pool. |
| Can parameters represent a table name? | Not as ordinary value parameters; use reviewed/whitelisted SQL structure. |
| Does BEGIN prevent every lost update? | No; statement shape, isolation, and locking still matter. |
| Does `EXPLAIN ANALYZE` only inspect? | No; it executes the statement. |
| Is replication a backup? | No; it can reproduce accidental corruption/deletion. |
| Does pipelining make several commands atomic? | No; it primarily reduces round trips. |
| Can outbox publication duplicate? | Yes; design consumers and operation identities accordingly. |
