# Context: USD Valuation of Binance Balances

Interactive change (no external spec).

## Existing code relevant to this change
- `modules/binance/utils.go` — `FormatAssetSymbolsToGetValue` builds `<ASSET>USDT` pairs for a
  hardcoded popular list (`modules/binance/constants.go`: BNB, ETH, BTC, FDUSD) and `<ASSET>BTC`
  for everything else; `FormatStringSliceToJSONStringArray` serializes the list.
- `modules/binance/services.go` — `GetPairValues` calls `/api/v3/ticker/price?symbols=[...]`;
  `GetSpotData` (`/api/v3/account`), `GetFundData` → `GetFundingData` (`/sapi/v1/asset/get-funding-asset`).
- `modules/binance/structs.go` — `FundingAssetResponse` already has `BtcValuation`.
- `modules/binance/controllers.go` — `Test()` handler exercises conversion with hardcoded balances.
- `db/models/binance.model.go` — `BinanceWallet.Usd string` column exists but is never populated.

## Known problems with the current attempt
- One invalid/delisted symbol (e.g. EDGBTC, MEETONEBTC) makes Binance reject the whole batch.
- `<ASSET>BTC` prices are never multiplied by BTCUSDT → not actually USD.
- Stablecoins and Simple Earn `LD*` assets are not handled.

## Discovered components
- Existing Go backend (Gin + GORM + Postgres), module `goGinServer`. No new components;
  no tech pack registered.
