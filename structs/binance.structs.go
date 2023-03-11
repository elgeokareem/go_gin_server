package binanceStructs

type WalletBinance struct {
	MakerCommission  int64         `json:"makerCommission"`
	TakerCommission  int64         `json:"takerCommission"`
	BuyerCommission  int64         `json:"buyerCommission"`
	SellerCommission int64         `json:"sellerCommission"`
	CanTrade         bool          `json:"canTrade"`
	CanWithdraw      bool          `json:"canWithdraw"`
	CanDeposit       bool          `json:"canDeposit"`
	Balances         []BalanceData `json:"balances"`
}

type BalanceData struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

type BalanceDataParsed struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}
