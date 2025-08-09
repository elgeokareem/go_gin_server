package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/samber/lo"
)

func GetFundingData(endpoint string) ([]BalanceData, error) {
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_API_SECRET")

	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	queryString := "timestamp=" + timestamp
	signature := hmac.New(sha256.New, []byte(secretKey))
	signature.Write([]byte(queryString))
	signatureString := hex.EncodeToString(signature.Sum(nil))
	queryString += "&signature=" + signatureString

	client := &http.Client{}
	req, _ := http.NewRequest("POST", endpoint, nil)
	req.Header.Add("X-MBX-APIKEY", apiKey)
	req.URL.RawQuery = queryString

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	strBody := string(body)
	fmt.Println(strBody)

	// Unmarshal the JSON response into the struct
	var fundingData []*FundingAssetResponse
	if err := json.Unmarshal(body, &fundingData); err != nil {
		return nil, err
	}

	// Map body response to balance data
	result := make([]BalanceData, len(fundingData))
	for i, item := range fundingData {
		result[i] = BalanceData{
			Asset:  item.Asset,
			Free:   item.Free,
			Locked: item.Locked,
		}
	}

	return result, err
}

func FormatAssetSymbolsToGetValue(items []BalanceData) string {
	formattedSymbols := make([]string, 0)
	for _, item := range items {
		assetName := item.Asset
		_, foundSymbol := lo.Find(PopularSymbolList, func(item string) bool {
			return item == assetName
		})

		if assetName == "" || assetName == "USDT" {
			continue
		}
		if foundSymbol {
			formattedSymbols = append(formattedSymbols, assetName+"USDT")
		}
		// If the symbol is not popular we fallback to BTC comparison
		if assetName != "BTC" && !foundSymbol {
			formattedSymbols = append(formattedSymbols, assetName+"BTC")
		}
	}

	result, _ := FormatStringSliceToJSONStringArray(formattedSymbols)
	return result
}

func FormatStringSliceToJSONStringArray(symbols []string) (string, error) {
	if symbols == nil {
		return "[]", nil
	}

	jsonData, err := json.Marshal(symbols)
	if err != nil {
		return "", fmt.Errorf("failed to marshal symbols to JSON: %w", err)
	}
	return string(jsonData), nil
}
