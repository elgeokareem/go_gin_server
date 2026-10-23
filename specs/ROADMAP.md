# Roadmap: billetera-backend

Recorded on 2026-10-23, when the project was picked up again after several months and SDD was set up.

## Vision

Users connect a **read-only** exchange API key (Binance first). The backend:

1. fetches their balances,
2. runs a cron job that periodically records balances and account movements,
3. keeps the full history, so the client can show balance over time,
4. lets the user add a **reason** to each expense or movement.

## Status of the existing code

| Area | Status | Notes |
|------|--------|-------|
| Register / login (bcrypt + JWT) | Done | `modules/auth`; one JWT unit test |
| Auth middleware | Done, not applied | Commented out on `/binance` routes (`modules/binance/routes.go`) |
| Postgres connection + auto-create DB | Done | `db/InitDBService.go` |
| Models `User`, `BinanceWallet` | Partial | Balances stored as `float64`; `Usd` column never filled in |
| Fetch spot balances | Done | `GET /api/v3/account` |
| Fetch funding balances | Done | `POST /sapi/v1/asset/get-funding-asset` |
| USD conversion | Partial / broken | One invalid symbol fails the whole batch; BTC prices never converted to USD |
| Save balances to DB | Broken | `InsertSpotDataInDb` inserts an empty row; its call is commented out |
| Cron job | Not started | `cron/get-binance-funds.go` is a stub; gocron not installed |
| Movements / transactions | Not started | No model; no `reason` field |
| Per-user Binance credentials | Not started | Keys come from `.env`; `User.BinanceToken` is one field (key + secret needed) |
| Multi-user | Not started | Each user would be a separate tenant (undecided) |
| Migrations | Inconsistent | GORM AutoMigrate **and** golang-migrate (empty `db/migrations/`, wrong relative paths in `cmd/migrate`) |

## Known bugs and tech debt

- `GetSpotHandler` / `GetFundHandler` keep going after an error (missing `return`).
- Spot balances are sorted by `Free` as a string (lexicographic).
- DB connection string, password included, is printed to stdout (`db/InitDBService.go`).
- `bcrypt.MinCost` for password hashing.
- Login sets a cookie, but the middleware only reads the `Authorization` header; CORS doesn't allow `Authorization`; cookie TTL (1h) ≠ JWT TTL (24h).
- JWT audience is the placeholder `"somebody_else"`.
- `middlewares/cors.go` is unused (CORS is configured in `routes/routes.go`).
- 21 MB compiled binary `goGinServer` is committed to git.

## Planned changes (proposed order)

1. **usd-valuation** (in progress): value spot and funding balances in USD using Binance's own `btcValuation` × the BTC/USD price, and save the results.
2. **bugfixes**: the known bugs above, auth on `/binance` routes, one migration system, `decimal` instead of `float64`.
3. **per-user-credentials**: store each user's read-only API key and secret (encrypted at rest); remove the `.env` keys.
4. **balance-snapshots-cron**: gocron job that records balances and their valuation every N minutes for every user.
5. **movements**: pull account movements (deposits, withdrawals, Binance Pay, P2P, Convert, trades), dedupe them, save them, and let users attach a **reason**.
6. **client API**: endpoints the frontend needs for history, movements and reasons.
7. **later**: other exchanges, multi-tenant model.
