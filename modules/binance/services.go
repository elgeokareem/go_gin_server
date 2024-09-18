package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
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
	params.Add("omitZeroBalances", "false")
	params.Add("timestamp", timestamp)

	// generate signature
	query := params.Encode()
	signature := hmac.New(sha256.New, []byte(secretKey))
	signature.Write([]byte(query))
	signatureString := fmt.Sprintf("%x", signature.Sum(nil))

	// create request
	url := binanceBaseURL + "/api/v3/account?" + query + "&signature=" + signatureString

	println(url)
	// url := binanceBaseURL + "/sapi/v1/capital/config/getall?" + query + "&signature=" + signatureString
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add API key to headers
	req.Header.Add("X-MBX-APIKEY", apiKey)

	// send request
	client := http.DefaultClient
	resp, err := client.Do(req)

	// Codigo para ver los headers del response y ver cuanto weight se ha usado
	// for name, values := range resp.Header {
	// 	// Loop over all values for the name.
	// 	for _, value := range values {
	// 		fmt.Printf("%s: %s\n", name, value)
	// 	}
	// }

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

	listOfAssets := accountData.Balances

	// Remove assets with 0 balance
	epsilon := 0.000001 // variable to make comparison
	for i := 0; i < len(listOfAssets); i++ {
		freeNumber, errFree := strconv.ParseFloat(listOfAssets[i].Free, 64)
		if errFree != nil {
			return nil, err
		}

		if freeNumber < epsilon {
			// remove item from list
			listOfAssets = append(listOfAssets[:i], listOfAssets[i+1:]...)
			i--
		}
	}

	sort.Slice(listOfAssets, func(i, j int) bool {
		return listOfAssets[i].Free > listOfAssets[j].Free
	})

	return listOfAssets, nil
}

func GetFundData() ([]*FundingAssetResponse, error) {
	endpoint := "https://api.binance.com/sapi/v1/asset/get-funding-asset"
	result, err := GetFundingData(endpoint)

	if err != nil {
		return nil, err
	}

	return result, nil
}
