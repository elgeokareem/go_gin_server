---
title: "Implementation Plan: USD Valuation of Binance Balances"
change: usd-valuation-1
type: feature
spec: ./SPEC.md
issue: none
created: 2026-10-23
sdd_version: 7.3.0
---

# Implementation Plan: USD Valuation of Binance Balances

## Overview

**Spec:** [SPEC.md](./SPEC.md)
**Issue:** none

There's no tech pack registered (Go codebase), so phases have no manifest agent or standards. Each phase is carried out by the main agent, following the conventions already in the repo:
- `gin.HandlerFunc` factories in `controllers.go`
- service functions in `services.go`
- GORM models in `db/models/`
- the global `db.Service.DB`

## Affected Components

- **backend** (existing, Go 1.27 / Gin / GORM / Postgres, module `goGinServer`): `modules/binance`, `db/models`, `db/initDatabases.go`, `go.mod`

No new components, so there's no scaffolding phase.

## Phases

### Phase 1: Binance client and valuation (pure logic)

**Outcome:** Signed requests for spot (`getUserAsset`) and funding (`get-funding-asset`) with `needBtcValuation=true`, the BTCUSDC price fetch, and decimal valuation functions, all unit-tested. No HTTP routes or DB changes yet.

**Deliverables:**
- `github.com/shopspring/decimal` added to `go.mod`
- One shared signing helper and request function, with the base URL from `BINANCE_BASE_URL` (removes the hardcoded host in the funding request)
- Response structs for user asset / funding (with `btcValuation`) and ticker
- Valuation: per-asset `usd_value` and `priced`, zero-amount skip, totals per wallet type and overall (SPEC FR4, Algorithm)
- Upstream error type carrying Binance `code` / `msg` (for the 502 body)
- Unit tests passing (see Tests: Unit)

**Spec refs:** FR1–FR4, NFR precision, NFR testability

### Phase 2: Persistence and sync service

**Outcome:** Snapshots can be saved atomically from Go code (no HTTP), ready for the future cron job to call.

**Deliverables:**
- `db/models/balance_snapshot.model.go`: `BalanceSnapshot` and `SnapshotAsset` with `numeric(38,18)` columns, FKs, check constraint, unique (`snapshot_id`, `wallet_type`, `asset`), index (`user_id`, `captured_at`) (SPEC Data Model)
- `db/initDatabases.go`: AutoMigrate both models (`binance_wallets` left alone)
- `SyncBalances(userID uint)`: fetch everything first, then one `db.Service.DB.Transaction`. Empty account still saves a snapshot.
- One structured log line per sync, without secrets (SPEC Observability)
- Persistence tests passing, gated by `TEST_DATABASE_URL`

**Spec refs:** FR6, NFR atomicity, AC2, AC5–AC7, Migration

**Depends on:** Phase 1

### Phase 3: HTTP layer, auth and cleanup

**Outcome:** The API contract in the SPEC is live behind JWT, and the old conversion code is gone.

**Deliverables:**
- `routes.go`: `AuthMiddleware` enabled on the `/binance` group; `POST /binance/sync` added; `/binance/test` removed
- `controllers.go`:
  - `GetSpotHandler` / `GetFundHandler` return the valued response shape
  - new `SyncHandler` returns 201
  - errors map to 401/502/500 with an early `return`
  - decimals serialized as strings
- Removed: `Test()`, `FormatAssetSymbolsToGetValue`, `FormatStringSliceToJSONStringArray`, `GetPairValues`, `PopularSymbolList`/`constants.go`, `InsertSpotDataInDb`, and the debug `fmt.Print*` calls in the module
- Handler tests against an `httptest.Server` fake Binance (handlers called with the user put in the gin context; the 401 test goes through the real middleware, which aborts before any DB access)

**Spec refs:** FR5, FR7, FR8, API Contract, Security, Error Handling, AC1, AC3, AC4

**Depends on:** Phase 2

### Phase 4: Integration and E2E testing

**Outcome:** Full suite green and a manual end-to-end check against the real account.

**Deliverables:**
- `go build ./...`, `go vet ./...`, `go test ./...` passing (persistence tests run with a local Postgres via `TEST_DATABASE_URL`)
- Manual E2E: `POST /auth/login` → `GET /binance/spot` → `POST /binance/sync` → check the `balance_snapshots` and `snapshot_assets` rows in Postgres
- Upstream call counts checked (3 per sync, 2 per live request)

**Depends on:** Phase 3

### Phase 5: Review

**Outcome:** Implementation checked against every SPEC acceptance criterion (AC1–AC7) and NFR; the roadmap entry updated; the user approves.

## Dependencies

- Binance API key from `.env` with read permission. A local Postgres for persistence tests and manual E2E.
- No dependency on other changes. Roadmap step 2 (bugfixes) builds on this one later.

## Tests

### Unit Tests
- [ ] `TestValuation_UsdIsBtcValuationTimesPrice`: exact decimal results
- [ ] `TestValuation_ZeroBtcValuationNonZeroAmount_IsUnpriced`: `usd_value` nil, `priced` false, left out of totals
- [ ] `TestValuation_ZeroAmountAssetsSkipped`
- [ ] `TestValuation_TotalsPerWalletTypeAndOverall`
- [ ] `TestValuation_SameAssetInSpotAndFundingKeptSeparate`
- [ ] `TestParse_UserAssetResponse`: string decimals, `btcValuation`
- [ ] `TestParse_FundingAssetResponse`
- [ ] `TestParse_TickerPrice`
- [ ] `TestSign_KnownHMAC`: fixed query + secret → expected signature
- [ ] `TestUpstreamError_CarriesBinanceCodeAndMsg`

### Integration Tests (handler, fake Binance)
- [ ] `TestSpotHandler_ReturnsValuedAssetsAndTotals` (AC1)
- [ ] `TestFundHandler_ReturnsValuedAssetsAndTotals` (AC1)
- [ ] `TestSpotHandler_UnpricedAssetSerializedAsNull` (AC3)
- [ ] `TestBinanceRoutes_NoJWT_Returns401_NoUpstreamCalls` (AC4)
- [ ] `TestSpotHandler_BinanceError_Returns502WithCode`
- [ ] `TestSpotHandler_MakesTwoUpstreamCalls`

### Persistence Tests (`TEST_DATABASE_URL`)
- [ ] `TestSyncBalances_PersistsSnapshotAndAssets` (AC2)
- [ ] `TestSyncBalances_PriceFailure_PersistsNothing` (AC5)
- [ ] `TestSyncBalances_EmptyAccount_SnapshotWithZeroTotal` (AC6)
- [ ] `TestSyncBalances_SameAssetTwoWallets_TwoRows` (AC7)
- [ ] `TestSyncHandler_Returns201WithSnapshot` (AC2)
- [ ] `TestSyncBalances_MakesThreeUpstreamCalls`

### E2E Tests
- [ ] Manual: login → spot → sync → inspect DB rows (Phase 4)

## Expected Files

### Files to Create
- `modules/binance/client.go`: signing, request, upstream error
- `modules/binance/valuation.go`
- `modules/binance/sync.go`: `SyncBalances`
- `modules/binance/valuation_test.go`, `client_test.go`, `controllers_test.go`, `sync_test.go`
- `db/models/balance_snapshot.model.go`

### Files to Modify
- `modules/binance/services.go`, `utils.go`, `structs.go`, `controllers.go`, `routes.go`
- `modules/binance/constants.go` (deleted)
- `db/initDatabases.go`
- `go.mod`, `go.sum`
- `specs/ROADMAP.md` (status of step 1)

## Risks

| Risk | Mitigation |
|------|------------|
| `getUserAsset` leaves out some spot assets, or `btcValuation` is 0 for tiny ("dust") balances | Treated as unpriced and flagged; checked during manual E2E against the real account |
| BTCUSDC price glitches / USDC moving away from $1 | Price and symbol stored with every snapshot, so values can be recomputed later |
| Persistence tests need Postgres | Skipped unless `TEST_DATABASE_URL` is set; pure logic covered by unit tests |
| AutoMigrate drift vs future golang-migrate | Accepted (SPEC A3); consolidated in roadmap step 2 |
| Binance rate limits during repeated manual testing | 3 calls per sync; tests use the fake server |

## Implementation State

- **Current Phase:** not started
- **Status:** pending
- **Completed Phases:**
  - [ ] Phase 1
  - [ ] Phase 2
  - [ ] Phase 3
  - [ ] Phase 4
  - [ ] Phase 5
- **Actual Files Changed:** none yet
- **Blockers:** none
- **Resource Usage:** n/a
