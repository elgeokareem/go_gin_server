package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"goGinServer/db"
	"goGinServer/db/models"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

func GetSpotData() ([]BalanceData, error) {
	binanceBaseURL := os.Getenv("BINANCE_BASE_URL")
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_API_SECRET")

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	params := url.Values{}
	params.Add("omitZeroBalances", "true")
	params.Add("timestamp", timestamp)

	// generate signature
	query := params.Encode()
	signature := hmac.New(sha256.New, []byte(secretKey))
	signature.Write([]byte(query))
	signatureString := fmt.Sprintf("%x", signature.Sum(nil))

	// create request
	url := binanceBaseURL + "/api/v3/account?" + query + "&signature=" + signatureString
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add API key to headers
	req.Header.Add("X-MBX-APIKEY", apiKey)

	// send request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response data as JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accountData Wallet
	err = json.Unmarshal(body, &accountData)
	if err != nil {
		return nil, err
	}

	balances := accountData.Balances

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Free > balances[j].Free
	})

	return balances, nil
}

func GetFundData() ([]BalanceData, error) {
	endpoint := "https://api.binance.com/sapi/v1/asset/get-funding-asset"
	result, err := GetFundingData(endpoint)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetPairValues(symbolList string) ([]TickerPriceResponse, error) {
	binanceBaseURL := os.Getenv("BINANCE_BASE_URL")

	escapedSymbolsValue := url.QueryEscape(symbolList)

	// create request
	fullURL := fmt.Sprintf("%s/api/v3/ticker/price?symbols=%s", binanceBaseURL, escapedSymbolsValue)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	// send request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response data as JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ticketPriceData []TickerPriceResponse
	err = json.Unmarshal(body, &ticketPriceData)
	if err != nil {
		return nil, err
	}

	return ticketPriceData, nil
}

func InsertSpotDataInDb(spotData []BalanceData, user models.User) error {
	spotRecord := models.BinanceWallet{UserID: user.ID}

	// Add user data to DB
	result := db.Service.DB.Create(&spotRecord)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
