---
title: USD Valuation of Binance Balances
spec_type: tech
type: feature
status: active
domain: Portfolio
issue: none
created: 2026-10-23
updated: 2026-10-23
sdd_version: 7.3.0
---

# USD Valuation of Binance Balances

## Overview

Value every Binance spot and funding balance in USD, using Binance's own `btcValuation` × the BTC/USDC price, and save the results as an append-only history of balance snapshots.

### Background

The product goal is to track a user's balance and expenses over time (see [ROADMAP.md](../ROADMAP.md)). Raw per-asset amounts don't give a single "what is this worth" number, and there's no history.

### Current State

- `modules/binance/utils.go` `FormatAssetSymbolsToGetValue` builds `<ASSET>USDT` pairs for a hardcoded list (BNB, ETH, BTC, FDUSD) and `<ASSET>BTC` for everything else. `GetPairValues` sends them all in one `/api/v3/ticker/price?symbols=[...]` request.
  - One invalid or delisted symbol makes Binance reject the whole batch.
  - `<ASSET>BTC` prices are never converted to USD.
  - Stablecoins and `LD*` assets are not handled.
- `GetSpotData` uses `/api/v3/account`, which returns no valuation.
- `GetFundingData` hardcodes `https://api.binance.com` and ignores `BINANCE_BASE_URL`.
- `BinanceWallet.Usd` (string) exists but is never filled in. `InsertSpotDataInDb` saves an empty row and its call is commented out.
- `/binance/*` routes have the auth middleware commented out.
- `GetSpotHandler` and `GetFundHandler` keep going after an error (missing `return`).

## User Stories

- **US1:** As a user, I want my live spot and funding balances to show their USD value, so I know what my account is worth right now.
- **US2:** As a user, I want to save a snapshot of my balances with their USD values, so I can see how my balance changes over time.
- **US3:** As a user, I want to see which assets couldn't be priced, so a missing value is never silently counted as $0.

## Functional Requirements

### FR1: Spot balances with BTC valuation

**Behavior:**
- The system shall fetch spot balances from `POST /sapi/v3/asset/getUserAsset` with `needBtcValuation=true` (signed request), replacing `/api/v3/account`.
- Each asset carries `free`, `locked`, `freeze` and `btcValuation`.

### FR2: Funding balances with BTC valuation

**Behavior:**
- The system shall call `POST /sapi/v1/asset/get-funding-asset` with `needBtcValuation=true`.
- The base URL comes from `BINANCE_BASE_URL` (no hardcoded host).

### FR3: BTC → USD price

**Behavior:**
- The system shall fetch `GET /api/v3/ticker/price?symbol=BTCUSDC` once per request or sync.
- USDC is treated as 1 USD.
- The price used is returned in responses and stored with each snapshot.

### FR4: Asset USD value and totals

**Behavior:**
- `usd_value = btc_valuation × btc_usd_price`, computed with exact decimals.
- When an asset's total amount (`free + locked + freeze`) > 0 and `btcValuation` = 0, the asset shall have `priced = false` and `usd_value = null`, and be left out of totals.
- Totals: `total_usd` per wallet type (`spot`, `funding`) and overall.

### FR5: Live endpoints

**Behavior:**
- `GET /binance/spot` and `GET /binance/fund` return live balances with valuation (FR4), wallet total and the BTC price used. Nothing is saved.

### FR6: Sync endpoint (save a snapshot)

**Behavior:**
- `POST /binance/sync` fetches spot, funding and the BTC price, computes the valuations, and saves one `balance_snapshots` row plus one `snapshot_assets` row per asset with a non-zero amount, all in one transaction.
- Snapshots are append-only (never updated or overwritten).
- The same asset in spot and funding is saved as two rows with different `wallet_type` values.
- An empty account still saves a snapshot with `total_usd = 0`.
- The sync logic is an exported service function that the future cron job (roadmap step 4) can call without HTTP.

### FR7: Authentication

**Behavior:**
- All `/binance/*` routes require a valid JWT (`middlewares.AuthMiddleware`). Snapshots belong to the authenticated user.
- Binance credentials still come from `BINANCE_API_KEY` / `BINANCE_API_SECRET` (per-user keys are roadmap step 3).

### FR8: Remove the old conversion code

**Behavior:**
- Remove `Test()` and its `/binance/test` route, `FormatAssetSymbolsToGetValue`, `FormatStringSliceToJSONStringArray`, `GetPairValues`, `PopularSymbolList`, and `InsertSpotDataInDb`.

## Non-Functional Requirements

| Requirement | Target | Measurement |
|-------------|--------|-------------|
| Monetary precision | No `float64` in the valuation path; `shopspring/decimal` in Go, `numeric` in Postgres | Code review; unit tests with exact values |
| Atomicity | A sync saves all of the snapshot or none of it | Persistence test with a failure injected |
| Binance API usage | 3 calls per sync / live request (spot or funding + price) | Handler tests count fake-server hits |
| Secret handling | API secret and signature never logged or returned | Code review; error responses checked in tests |
| Testability | Binance base URL injectable via `BINANCE_BASE_URL` | Handler tests run against `httptest.Server` |

## Technical Design

### Architecture

```
GET  /binance/spot ─┐
GET  /binance/fund ─┼─ AuthMiddleware ─ controllers ─ services ─┬─ binance client (signed requests, base URL from env)
POST /binance/sync ─┘                                             ├─ valuation (pure functions, decimal)
                                                                  └─ snapshot repository (GORM, one transaction)
```

- **Binance client** (`modules/binance`): one shared signing helper (HMAC-SHA256 over the query string, `X-MBX-APIKEY` header) used by spot and funding. It replaces the duplicated signing code in `GetSpotData` and `GetFundingData`.
- **Valuation:** pure functions from parsed balances and the BTC price to valued assets and totals, so they're unit-testable without HTTP or a DB.
- **Sync service:** `SyncBalances(userID uint) (BalanceSnapshot, error)`. It fetches everything first, then saves in one `db.Service.DB.Transaction`.

### Data Model

New GORM models in `db/models/`, created by GORM AutoMigrate through `db.InitDatabases`.

**`balance_snapshots`**

| Column | Type | Notes |
|--------|------|-------|
| id | bigserial PK | |
| user_id | bigint FK → users.id, not null, indexed | |
| source | varchar(16), not null | `binance` |
| btc_usd_price | numeric(38,18), not null | BTCUSDC price used |
| price_symbol | varchar(16), not null | `BTCUSDC` |
| total_usd | numeric(38,18), not null | priced assets only |
| spot_total_usd | numeric(38,18), not null | |
| funding_total_usd | numeric(38,18), not null | |
| unpriced_count | int, not null | |
| captured_at | timestamptz, not null, indexed with user_id | |
| created_at | timestamptz | |

**`snapshot_assets`**

| Column | Type | Notes |
|--------|------|-------|
| id | bigserial PK | |
| snapshot_id | bigint FK → balance_snapshots.id, not null, indexed, on delete cascade | |
| wallet_type | varchar(16), not null, check in (`spot`,`funding`) | |
| asset | varchar(32), not null | |
| free | numeric(38,18), not null | |
| locked | numeric(38,18), not null | |
| freeze | numeric(38,18), not null | |
| btc_valuation | numeric(38,18), not null | as returned by Binance |
| usd_value | numeric(38,18), null | null when unpriced |
| priced | boolean, not null | |

Unique constraint on (`snapshot_id`, `wallet_type`, `asset`).

The existing `binance_wallets` table and `BinanceWallet` model are left alone in this change (removed in the bugfixes change).

### Algorithm

```
price  = GET ticker/price BTCUSDC
spot   = getUserAsset(needBtcValuation=true)
fund   = get-funding-asset(needBtcValuation=true)
for each asset in spot ∪ fund (tagged by wallet type):
    amount = free + locked + freeze
    skip if amount == 0
    if btcValuation > 0: usd = btcValuation × price; priced = true
    else:                usd = null;                 priced = false
totals = sum(usd where priced), per wallet type and overall
```

## API Contract

All routes need `Authorization: Bearer <jwt>`. Decimal values are serialized as **strings**.

### `GET /binance/spot` and `GET /binance/fund` → 200

```json
{
  "wallet_type": "spot",
  "price_symbol": "BTCUSDC",
  "btc_usd_price": "67250.10000000",
  "total_usd": "2154.31",
  "unpriced_count": 1,
  "assets": [
    { "asset": "BTC", "free": "0.00581088", "locked": "0", "freeze": "0",
      "btc_valuation": "0.00581088", "usd_value": "390.78", "priced": true },
    { "asset": "MEETONE", "free": "0.8991", "locked": "0", "freeze": "0",
      "btc_valuation": "0", "usd_value": null, "priced": false }
  ]
}
```

### `POST /binance/sync` → 201

```json
{
  "snapshot": {
    "id": 12,
    "captured_at": "2026-10-23T14:00:00Z",
    "price_symbol": "BTCUSDC",
    "btc_usd_price": "67250.10000000",
    "total_usd": "2201.90",
    "spot_total_usd": "2154.31",
    "funding_total_usd": "47.59",
    "unpriced_count": 1,
    "assets": [ { "wallet_type": "spot", "asset": "BTC", "...": "same fields as above" } ]
  }
}
```

### Errors

| Status | When | Body |
|--------|------|------|
| 401 | Missing, malformed, expired or invalid JWT | `{"error": "..."}` (existing middleware messages) |
| 502 | Binance or price request failed / non-2xx / unparsable | `{"error": "binance request failed", "binance_code": -2015, "binance_msg": "..."}` (code and msg only when Binance returned them) |
| 500 | DB transaction failed | `{"error": "could not save snapshot"}` |

## Security Considerations

- `/binance/*` is protected by `AuthMiddleware` (FR7).
- Binance API secret, signature and full signed URL are never logged or included in error responses. The existing `fmt.Println(strBody)` / `fmt.Printf("Data: ...")` debug prints in this module are removed.
- Snapshots are always scoped to the authenticated `user_id`.

## Error Handling

- Any upstream failure (network, non-2xx, JSON decode) during sync aborts **before** anything is saved → 502. Nothing is written.
- DB failure inside the transaction → rollback → 500.
- Handlers `return` right after writing an error response.

## Observability

- One log line per sync: user id, snapshot id, asset count, unpriced count, total USD, duration. No secrets.
- On upstream failure: log the endpoint path (no query string) and the Binance error code.

## Acceptance Criteria

- **AC1:** Given valid Binance keys and an authenticated user, when they call `GET /binance/spot` or `GET /binance/fund`, then each asset includes `btc_valuation`, `usd_value` and `priced`, and the response includes `total_usd`, `btc_usd_price` and `price_symbol`.
- **AC2:** Given valid Binance keys and an authenticated user, when they call `POST /binance/sync`, then one `balance_snapshots` row and one `snapshot_assets` row per non-zero asset (spot and funding) are saved for that user, and the saved snapshot comes back with 201.
- **AC3:** Given an asset with a non-zero amount and `btcValuation` = 0, when it's valued, then it has `priced = false` and `usd_value = null`, and is left out of all totals.
- **AC4:** Given a request without a valid JWT, when it calls any `/binance/*` route, then the response is 401 and no Binance request is made.
- **AC5:** Given Binance or the BTCUSDC price request fails, when the user calls `POST /binance/sync`, then the response is 502 and no snapshot or asset rows are saved.
- **AC6:** Given an account with no non-zero balances, when the user calls `POST /binance/sync`, then a snapshot with `total_usd = 0` and no asset rows is saved.
- **AC7:** Given the same asset in both spot and funding, when a sync runs, then two `snapshot_assets` rows are saved with different `wallet_type` values.

## Domain Model

### Entities

| Entity | Definition | Spec Path | Status |
|--------|------------|-----------|--------|
| User | Account holder authenticated via JWT | existing | existing |
| Balance Snapshot | Point-in-time record of a user's Binance holdings and their USD value | this spec | new |
| Snapshot Asset | One asset's amount and valuation within a snapshot, for one wallet type | this spec | new |

### Relationships

```
User 1 ──── * BalanceSnapshot 1 ──── * SnapshotAsset
```

### Glossary

| Term | Definition | First Defined In |
|------|------------|------------------|
| Wallet type | Binance wallet a balance belongs to: `spot` or `funding` | this spec |
| BTC valuation | Binance-computed value of an asset holding, in BTC (`btcValuation`) | this spec |
| BTC/USD price | Last BTCUSDC price, with USDC treated as 1 USD | this spec |
| Unpriced asset | Asset with a non-zero amount but zero BTC valuation | this spec |
| Sync | Fetching, valuing and saving a balance snapshot | this spec |

### Bounded Contexts

- **Portfolio:** balances, valuations, snapshots (this change). Later: movements and reasons.
- **Identity:** users and auth (existing).

## Specs Directory Changes

Layout convention: each change lives in `specs/YYYY-MM-DD-<name>/` with its `SPEC.md` and `PLAN.md`.

```
Before                         After
specs/                         specs/
└── INDEX.md                   ├── INDEX.md
                               ├── ROADMAP.md
                               └── 2026-10-23-usd-valuation/
                                   ├── SPEC.md
                                   └── PLAN.md
```

| Path | Action | Description |
|------|--------|-------------|
| specs/INDEX.md | Modified | Lists this change |
| specs/ROADMAP.md | Added | Project vision, current status, planned changes |
| specs/2026-10-23-usd-valuation/SPEC.md | Added | This spec |
| specs/2026-10-23-usd-valuation/PLAN.md | Added | Implementation plan |

## Components

New components will be scaffolded during implementation. None are needed here.

### New Components

| Component | Type | Settings | Purpose |
|-----------|------|----------|---------|
| *none* | | | |

### Modified Components

| Component | Changes |
|-----------|---------|
| backend (Go, `goGinServer`, no tech pack) | `modules/binance` (client, valuation, sync, routes), `db/models` (2 new models), `db/initDatabases.go` (AutoMigrate), `go.mod` (`shopspring/decimal`) |

## System Analysis

### Inferred Requirements

- A shared signing helper is needed for both signed endpoints (inferred from duplicated code).
- The funding request must use `BINANCE_BASE_URL` so it can be tested (inferred from NFR testability).

### Gaps & Assumptions

- **A1:** Binance returns `btcValuation = "0"` for assets it can't value. Treated as unpriced (FR4).
- **A2:** USDC ≈ 1 USD is accurate enough (user-approved, option a).
- **A3:** GORM AutoMigrate is used for the new tables until migrations are consolidated (roadmap step 2).

### Dependencies

- Binance endpoints: `/sapi/v3/asset/getUserAsset`, `/sapi/v1/asset/get-funding-asset`, `/api/v3/ticker/price`.
- The existing `middlewares.AuthMiddleware` and `auth.GetUserFromContext`.

## Requirements Discovery

### Component Discovery Phase

| # | Question | Answer | Source |
|---|----------|--------|--------|
| D1 | Which components are involved? | Existing Go backend only (`modules/binance`); no new components, no tech pack | Assumption (confirmed by context) |

### Solicitation Phase

| # | Question | Answer | Source |
|---|----------|--------|--------|
| S1 | Confirm problem / primary user / expected outcome draft | Confirmed; goal is to track balance and all expenses over time | User |
| S2 | Pricing source: (a) full ticker map, (b) Binance `btcValuation` × BTC price, or a with b as backup? | Option b. If Binance already gives the answer, no need to recreate it | User |
| S3 | Is treating USDT as USD acceptable? | USD is preferable | User |
| S4 | Valuations only in API responses, or saved? | Saved. The idea is to track the balance and all expenses in general | User |
| S5 | Anything to add or remove? | Add a general roadmap recording everything discussed | User |
| S6 | True USD: (a) BTCUSDC with USDC = $1, or (b) BTCUSDT corrected by an external USDT/USD rate? | (a) | User |
| S7 | Append-only `balance_snapshots` + `snapshot_assets`, decimal `numeric`, replacing `BinanceWallet.Usd`? | Yes | User |
| S8 | Snapshot trigger: dedicated `POST /binance/sync`, or save on every GET? | Sync endpoint | User |
| S9 | Re-enable auth on `/binance` routes in this change; keys still from `.env`? | Yes | User |
| S10 | Earn / staking / futures / margin out of scope, add to roadmap? | Agree | User |
| S11 | Confirm drafted NFRs, user stories, ACs, edge cases, dependencies, tests, out of scope | Ok | User |

### User Feedback & Corrections

- 2026-10-23: Asked for the roadmap of all missing and partially done work to be recorded → `specs/ROADMAP.md`.

### Open Questions (BLOCKING)

| # | Question | Status | Blocker For |
|---|----------|--------|-------------|
| *none* | | | |

## Testing Strategy

### Unit Tests

| Test | Covers |
|------|--------|
| `usd_value = btc_valuation × price` with exact decimal results | FR4, AC1 |
| Zero `btcValuation` with non-zero amount → unpriced, left out of totals | FR4, AC3 |
| Zero-amount assets skipped | FR6 |
| Totals per wallet type and overall | FR4 |
| Parsing `getUserAsset`, funding and ticker JSON (string decimals) | FR1–FR3 |
| Signing helper output matches a known HMAC for a fixed query | FR1, FR2 |

### Integration Tests (handler level, `httptest` fake Binance)

| Test | Covers |
|------|--------|
| `GET /binance/spot` and `/fund` succeed with valuation fields | AC1 |
| No JWT → 401 and the fake Binance server gets 0 hits | AC4 |
| Binance returns `{"code":-2015,...}` → 502 with `binance_code` | Error handling |
| Exactly 2 upstream calls per live request, 3 per sync | NFR API usage |

### Persistence Tests (skipped unless `TEST_DATABASE_URL` is set)

| Test | Covers |
|------|--------|
| Sync saves the snapshot plus the expected asset rows for the user | AC2 |
| Price request fails → no rows saved | AC5 |
| Empty account → snapshot with `total_usd = 0`, no assets | AC6 |
| Same asset in spot and funding → 2 rows | AC7 |

### E2E Tests

Not automated. Manual check against the real Binance account: `POST /auth/login` → `POST /binance/sync` → inspect the rows in Postgres.

### Test Data

Fixtures based on the balances in the current `Test()` handler (EDG, FLR, USDT, BNB, ETH, BTC, MEETONE …), including one unpriced asset.

## Dependencies

- **Internal:** `db.Service`, `middlewares.AuthMiddleware`, `auth.GetUserFromContext`, `models.User`.
- **External:** Binance REST API (same API key, read permission enough); Go module `github.com/shopspring/decimal`.

## Migration / Rollback

- **Forward:** AutoMigrate creates `balance_snapshots` and `snapshot_assets` at startup. No existing data is changed.
- **Rollback:** revert the code; drop the two new tables (`DROP TABLE snapshot_assets, balance_snapshots;`). `binance_wallets` is untouched.

## Out of Scope

- Simple Earn, staking, futures and margin balances (added to the roadmap)
- The cron job / scheduled syncs (roadmap step 4)
- Per-user Binance credentials (roadmap step 3)
- Movements, expenses and reasons (roadmap step 5)
- Removing `binance_wallets` and consolidating migrations (roadmap step 2)

## Open Questions

None.

## References

- [specs/ROADMAP.md](../ROADMAP.md)
- Binance: User Asset `POST /sapi/v3/asset/getUserAsset`; Funding Wallet `POST /sapi/v1/asset/get-funding-asset`; Symbol Price Ticker `GET /api/v3/ticker/price`
