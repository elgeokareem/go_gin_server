package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	structsBinance "goGinServer/structs"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

func GetWalletData() []structsBinance.BalanceData {
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_API_SECRET")

	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	params := url.Values{}
	params.Add("timestamp", timestamp)

	// generate signature
	query := params.Encode()
	signature := hmac.New(sha256.New, []byte(secretKey))
	signature.Write([]byte(query))
	signatureString := fmt.Sprintf("%x", signature.Sum(nil))

	// create request
	url := "https://api.binance.com/api/v3/account?" + query + "&signature=" + signatureString
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	// add API key to headers
	req.Header.Add("X-MBX-APIKEY", apiKey)

	// send request
	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// parse response data as JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var accountData structsBinance.Wallet
	err = json.Unmarshal(body, &accountData)
	if err != nil {
		panic(err)
	}

	listOfAssets := accountData.Balances

	// Remove assets with 0 balance
	epsilon := 0.000001 // variable to make comparison
	for i := 0; i < len(listOfAssets); i++ {
		freeNumber, errFree := strconv.ParseFloat(listOfAssets[i].Free, 64)
		if errFree != nil {
			panic(errFree)
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

	return listOfAssets
}
