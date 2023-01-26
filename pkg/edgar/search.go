package edgar

import (
	"net/http"
	"time"
)

// TODO: convert date and datetime fields to time.Time

// Company is a simple struct for a single company.
type Company struct {
	Symbol              *string `json:"symbol"`
	Exchange            *string `json:"exchange"`
	ExchangeSuffix      *string `json:"exchangeSuffix"`
	ExchangeName        *string `json:"exchangeName"`
	ExchangeSegment     *string `json:"exchangeSegment"`
	ExchangeSegmentName *string `json:"exchangeSegmentName"`
	Name                *string `json:"name"`
	Date                *string `json:"date"` // TODO: make time.Time
	Type                *string `json:"type"`
	IexID               *string `json:"iexId"`
	Region              *string `json:"region"`
	Currency            *string `json:"currency"`
	IsEnabled           *bool   `json:"isEnabled"`
	FIGI                *string `json:"figi"`
	CIK                 *string `json:"cik"`
	LEI                 *string `json:"lei"`
}

// SubmissionsSearchResult is a single sec filer.
type SubmissionsSearchResult struct {
	CIK                               *string       `json:"cik"`
	EntityType                        *string       `json:"entityType"`
	SIC                               *string       `json:"sic"`
	SICDescription                    *string       `json:"sicDescription"`
	InsiderTransactionForOwnerExists  *int8         `json:"insiderTransactionForOwnerExists"`
	InsiderTransactionForIssuerExists *int8         `json:"insiderTransactionForIssuerExists"`
	Name                              *string       `json:"name"`
	Tickers                           *[]string     `json:"tickers"`
	Exchanges                         *[]string     `json:"exchanges"`
	EIN                               *string       `json:"ein"`
	Description                       *string       `json:"description"`
	Website                           *string       `json:"website"`
	InvestorWebsite                   *string       `json:"investorWebsite"`
	Category                          *string       `json:"category"`
	FiscalYearEnd                     *string       `json:"fiscalYearEnd"`
	StateOfIncorporation              *string       `json:"stateOfIncorporation"`
	StateOfIncorporationDescription   *string       `json:"stateOfIncorporationDescription"`
	Addresses                         *Addresses    `json:"addresses"`
	Phone                             *string       `json:"phone"`
	Flags                             *string       `json:"flags"`
	FormerNames                       *[]NameRecord `json:"formerNames"`
	Filings                           *FilingRecord `json:"filings"`
}

type Addresses struct {
	Mailing  *Address `json:"mailing"`
	Business *Address `json:"business"`
}

type Address struct {
	Street1                   string `json:"street1"`
	Street2                   string `json:"street2"`
	City                      string `json:"city"`
	StateOrCountry            string `json:"stateOrCountry"`
	ZipCode                   string `json:"zipCode"`
	StateOrCountryDescription string `json:"stateOrCountryDescription"`
}

type NameRecord struct {
	Name string    `json:"name"`
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type FilingRecord struct {
	Recent *FilingTable `json:"recent"`
	Files  *[]string    `json:"file"`
}

type FilingTable struct {
	AccessionNumber       *[]string    `json:"accessionNumber"`
	FilingDate            *[]string    `json:"filingDate"`
	ReportDate            *[]string    `json:"reportDate"`
	AcceptanceDateTime    *[]time.Time `json:"acceptanceDateTime"`
	Act                   *[]string    `json:"act"`
	Form                  *[]string    `json:"form"`
	FileNumber            *[]string    `json:"fileNumber"`
	FilmNumber            *[]string    `json:"filmNumber"`
	Items                 *[]string    `json:"items"`
	Size                  *[]int64     `json:"size"`
	IsXBRL                *[]int8      `json:"isXBRL"`
	IsInlineXBRL          *[]int8      `json:"isInlineXBRL"`
	PrimaryDocument       *[]string    `json:"primaryDocument"`
	PrimaryDocDescription *[]string    `json:"primaryDocDescription"`
}

// TODO: in client.go, convert FilingTable results to []Filing (column major -> row major)
type Filing struct {
	CIK                   *string    `json:"cik"` // NOTE: this is a hack for easy unpacking of values in steampipe
	AccessionNumber       *string    `json:"accessionNumber"`
	FilingDate            *string    `json:"filingDate"`
	ReportDate            *string    `json:"reportDate"`
	AcceptanceDateTime    *time.Time `json:"acceptanceDateTime"`
	Act                   *string    `json:"act"`
	Form                  *string    `json:"form"`
	FileNumber            *string    `json:"fileNumber"`
	FilmNumber            *string    `json:"filmNumber"`
	Items                 *string    `json:"items"`
	Size                  *int64     `json:"size"`
	IsXBRL                *int8      `json:"isXBRL"`
	IsInlineXBRL          *int8      `json:"isInlineXBRL"`
	PrimaryDocument       *string    `json:"primaryDocument"`
	PrimaryDocDescription *string    `json:"primaryDocDescription"`
}

// Search methods
// ------------------

// GetPublicCompanies returns a list of public companies.
func (c *client) GetPublicCompanies() (*[]Company, error) {
	out := new([]Company)

	resp, err := c.request(http.MethodGet, iexSymbolsURL, nil)
	if err != nil {
		return out, err
	}

	err = unmarshall(resp, out)
	return out, err
}

// Filings methods
// -----------------

// TODO: handle pagination in the SEC Edgar API
// GetFilings gets a list of filings for a single CIK.
func (c *client) GetSubmissions(cik string) (submissions *SubmissionsSearchResult, err error) {
	submissions = new(SubmissionsSearchResult)

	url := secCompanyURL + "CIK" + cik + ".json"
	// url := "https://data.sec.gov/submissions/CIK0001650373.json"
	resp, err := c.request(http.MethodGet, url, nil)

	if err != nil {
		return submissions, err
	}

	err = unmarshall(resp, submissions)
	return submissions, err
}
