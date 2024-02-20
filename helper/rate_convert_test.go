package helper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRate_Convert(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
            "meta": {
                "last_updated_at": "2022-01-01T00:00:00Z"
            },
            "data": {
                "EUR": {
                    "code": "EUR",
                    "value": 1.2
                }
            }
        }`)
	}))
	defer server.Close()

	// Temporarily replace the API_KEY with the mock server's URL for testing
	originalAPIKey := API_KEY
	API_KEY = server.URL // This is just for illustrative purposes; adjust as needed for your actual code structure
	defer func() { API_KEY = originalAPIKey }()

	// Test cases
	tests := []struct {
		name           string
		baseCurrency   string
		targetCurrency string
		baseAmount     float64
		want           float64
		expectErr      bool
	}{
		{
			name:           "Successful conversion from USD to EUR",
			baseCurrency:   "USD",
			targetCurrency: "EUR",
			baseAmount:     100,
			want:           120, // Expected result based on the mock API response
			expectErr:      false,
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Rate_Convert_URL(server.URL, tt.baseAmount)
			if (err != nil) != tt.expectErr {
				t.Errorf("Rate_Convert() error = %v, wantErr %v", err, tt.expectErr)
			}
			if got != tt.want {
				t.Errorf("Rate_Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
