package edgar

// TODO: refactor this code and create a client library similar to https://github.com/calligram-engineering/steampipe-plugin-crunchbase/blob/056f7b1f3666fdbb460fb29a14ff511468643617/pkg/crunchbase/crunchbase.go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO: determine how to break dependency on IEX as it costs $49.99/month
//  1. Can just use IEX to get CIK by name and pay the cost
//  2. Can use EDGAR `curl -d "company=Google" -X POST https://www.sec.gov/cgi-bin/cik_lookup`
//     a. Need to implement some sort of keystroke search and fuzzy string search for when customers are searching for a company
const (
	iexSymbolsURL = "https://cloud.iexapis.com/stable/ref-data/symbols"
	secCompanyURL = "https://data.sec.gov/submissions/"
)

// Client definition
// ---------------------

type Client interface {
	GetPublicCompanies() (*[]Company, error)
	GetSubmissions(cik string) (*SubmissionsSearchResult, error) // TODO: add time window function
}

type client struct {
	iexToken string
	headers  map[string]string
}

// NewClient returns a pointer to a new EDGR Piquette client
func NewClient(token string) *client {
	c := client{}
	c.iexToken = token
	return &c
}

func (c *client) request(method, url string, body interface{}) (*http.Response, error) {
	payload, err := marshall(body)
	if err != nil {
		return nil, err
	}

	return c.do(method, url, bytes.NewReader(payload))
}

func (c *client) do(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// if we are using the IEX API, add the token to the request parameters
	if strings.Contains(url, "cloud.iexapis.com") {
		// appending to existing query args
		q := req.URL.Query()
		q.Add("token", c.iexToken)

		// assign encoded query string to http request
		req.URL.RawQuery = q.Encode()
	}

	c.headers = make(map[string]string)
	if body != nil {
		c.headers["Content-Type"] = "application/json"
	}
	// NOTE: see https://stackoverflow.com/questions/68131406/downloading-files-from-sec-gov-via-edgar-using-python-3-9
	// TODO: replace with headers on https://www.sec.gov/os/webmaster-faq#developers
	c.headers["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"
	c.headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func marshall(in interface{}) ([]byte, error) {
	if in == nil {
		return nil, nil
	}

	return json.Marshal(in)
}

func unmarshall(res *http.Response, out interface{}) error {
	fmt.Println(res)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		apiErr := new(APIError)
		apiErr.Response = res
		if err := json.NewDecoder(res.Body).Decode(apiErr); err != nil {
			apiErr.Value = Ptr("Oops! Something went wrong when parsing the error response.")
		}
		return apiErr
	}

	if out != nil {
		return json.NewDecoder(res.Body).Decode(out)
	}

	return nil
}

// --- API Error Responses ---

// APIError represents an error response returnted by the API.
type APIError struct {
	Response *http.Response

	Value *string `json:"value,omitmepty"`
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("%v %v: %d", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode)

	if e.Value != nil {
		msg = fmt.Sprintf("%s %v", msg, *e.Value)
	}

	return msg
}

func Ptr[T any](v T) *T {
	return &v
}
