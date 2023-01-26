package edgar

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetPublicCompanies calls GetPublicCompanies with a basic *[]Companies
// for a valid return value.
func TestGetPublicCompanies(t *testing.T) {
	apiKey := os.Getenv("IEX_API_KEY")
	fmt.Println("IEX API key = ", apiKey)
	if apiKey == "" {
		t.Fatalf("No IEX API Key found")
	}
	client := NewClient(apiKey)
	result, err := client.GetPublicCompanies()
	require.NoError(t, err)
	for i, org := range *result {
		if org.ExchangeName != nil {
			fmt.Println("Company ", i, " = ", *org.ExchangeName)
		}
	}
	// jsonBytes, err := json.Marshal(result)
	// fmt.Println(string(jsonBytes))
}

// TestGetSubmissions calls GetSubmissions with a basic *SubmissionsSearchResult
// for a valid return value.
func TestGetSubmissions(t *testing.T) {
	apiKey := os.Getenv("IEX_API_KEY")
	fmt.Println("IEX API key = ", apiKey)
	if apiKey == "" {
		t.Fatalf("No IEX API Key found")
	}
	client := NewClient(apiKey)
	result, err := client.GetSubmissions("0001650373")
	require.NoError(t, err)
	for i, num := range *result.Filings.Recent.IsInlineXBRL {
		fmt.Println("Company ", i, " = ", num)
	}
	filing := new(Filing)
	filing.CIK = result.CIK
	fmt.Println(*filing.CIK)
	// jsonBytes, err := json.Marshal(result)
	// fmt.Println(string(jsonBytes))
}
