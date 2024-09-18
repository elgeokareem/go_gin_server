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
)

func GetFundingData(endpoint string) ([]*FundingAssetResponse, error) {
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
	var result []*FundingAssetResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, err
}
