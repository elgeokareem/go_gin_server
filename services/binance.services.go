package services

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

func GetWalletData() {
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

	var accountData any
	err = json.Unmarshal(body, &accountData)
	if err != nil {
		panic(err)
	}

	listOfAssets := accountData.Balances

	// Make new slice
	// parsedData := make([]BinanceStructs.BalanceDataParsed, 0)

	// Remove assets with 0 balance
	for i := 0; i < len(listOfAssets); i++ {
		freeNumber, errFree := strconv.ParseFloat(listOfAssets[i].Free, 64)
		if errFree != nil {
			panic(errFree)
		}

		// lockedNumber, errLocked := strconv.ParseFloat(listOfAssets[i].Free, 64)
		// if errLocked != nil {
		// 	panic(errLocked)
		// }

		// parsedData = append(parsedData)

		if freeNumber == 0 {
			listOfAssets = append(listOfAssets[:i], listOfAssets[i+1:]...)
		}
	}

	sort.Slice(listOfAssets, func(i, j int) bool {
		return listOfAssets[i].Free > listOfAssets[j].Free
	})

	// for _, balance := range listOfAssets {
	// 	fmt.Printf("  %s: %s (locked: %s)\n", balance.Asset, balance.Free, balance.Locked)
	// }
}
