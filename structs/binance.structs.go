package structsBinance

type Wallet struct {
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

type GetAllCoins []struct {
	Coin             string `json:"coin"`
	DepositAllEnable bool   `json:"depositAllEnable"`
	Free             string `json:"free"`
	Freeze           string `json:"freeze"`
	Ipoable          string `json:"ipoable"`
	Ipoing           string `json:"ipoing"`
	IsLegalMoney     bool   `json:"isLegalMoney"`
	Locked           string `json:"locked"`
	Name             string `json:"name"`
	NetworkList      []struct {
		AddressRegex            string `json:"addressRegex"`
		Coin                    string `json:"coin"`
		DepositDesc             string `json:"depositDesc,omitempty"`
		DepositEnable           bool   `json:"depositEnable"`
		IsDefault               bool   `json:"isDefault"`
		MemoRegex               string `json:"memoRegex"`
		MinConfirm              int    `json:"minConfirm"`
		Name                    string `json:"name"`
		Network                 string `json:"network"`
		ResetAddressStatus      bool   `json:"resetAddressStatus"`
		SpecialTips             string `json:"specialTips"`
		UnLockConfirm           int    `json:"unLockConfirm"`
		WithdrawDesc            string `json:"withdrawDesc,omitempty"`
		WithdrawEnable          bool   `json:"withdrawEnable"`
		WithdrawFee             string `json:"withdrawFee"`
		WithdrawIntegerMultiple string `json:"withdrawIntegerMultiple"`
		WithdrawMax             string `json:"withdrawMax"`
		WithdrawMin             string `json:"withdrawMin"`
		SameAddress             bool   `json:"sameAddress"`
		EstimatedArrivalTime    int    `json:"estimatedArrivalTime"`
		Busy                    bool   `json:"busy"`
	} `json:"networkList"`
	Storage           string `json:"storage"`
	Trading           bool   `json:"trading"`
	WithdrawAllEnable bool   `json:"withdrawAllEnable"`
	Withdrawing       string `json:"withdrawing"`
}
