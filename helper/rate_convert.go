package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var API_KEY string

type ApiResponse struct {
	Meta MetaData            `json:"meta"`
	Data map[string]Currency `json:"data"`
}

type MetaData struct {
	LastUpdatedAt string `json:"last_updated_at"`
}

type Currency struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

func Rate_Convert(baseCurrency, targetCurrency string, base_amount float64) (float64, error) {

	if base_amount == 0 {
		return 0, nil
	}

	if !(len(API_KEY) > 0) {
		return 0, errors.New("CURRENCY_EXCHANGE_API_KEY missing")
	}

	url := fmt.Sprintf("https://api.currencyapi.com/v3/latest?apikey=%s&base_currency=%s&currencies=%s", API_KEY, baseCurrency, targetCurrency)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Decode JSON into the ApiResponse struct
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}

	var updated_amount float64

	fmt.Printf("Last Updated: %s\n", apiResponse.Meta.LastUpdatedAt)
	for code, currency := range apiResponse.Data {
		fmt.Printf("Currency: %s, Value: %f\n", code, currency.Value)
		updated_amount = currency.Value * base_amount
	}

	return updated_amount, nil
}

func Rate_Convert_URL(serverURL string, base_amount float64) (float64, error) {

	if base_amount == 0 {
		return 0, nil
	}

	url := serverURL

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Decode JSON into the ApiResponse struct
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Fatalf("Failed to decode JSON response: %v", err)
	}

	var updated_amount float64

	fmt.Printf("Last Updated: %s\n", apiResponse.Meta.LastUpdatedAt)
	for code, currency := range apiResponse.Data {
		fmt.Printf("Currency: %s, Value: %f\n", code, currency.Value)
		updated_amount = currency.Value * base_amount
	}

	return updated_amount, nil
}
