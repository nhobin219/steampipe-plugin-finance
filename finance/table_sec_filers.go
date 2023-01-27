package finance

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-finance/pkg/edgar"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableSecFilers(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "sec_filers",
		Description: "Lookup company filer details from the US SEC Edgar database.",
		List: &plugin.ListConfig{
			Hydrate:    listSecFiler,
			KeyColumns: plugin.SingleColumn("cik"),
		},
		Columns: []*plugin.Column{
			{Name: "cik", Type: proto.ColumnType_STRING, Transform: transform.FromField("CIK").Transform(transformCIK), Description: "CIK (Central Index Key) of the filer."},
			{Name: "entity_type", Type: proto.ColumnType_STRING, Description: "Entity type of the filer."},
			{Name: "sic", Type: proto.ColumnType_INT, Transform: transform.FromField("SIC").Transform(transformStringToInt), Description: "SIC (Standard Industrial Classification) of the filer."},
			{Name: "sic_description", Type: proto.ColumnType_STRING, Transform: transform.FromField("SICDescription"), Description: "SIC (Standard Industrial Classification) description of the filer."},
			{Name: "insider_transaction_for_owner_exists", Type: proto.ColumnType_INT, Description: "Whether or not an insider transaction for ther issuer of the filer exists."},
			{Name: "insider_transaction_for_issuer_exists", Type: proto.ColumnType_INT, Description: "Whether or not an insider transaction for ther owner of the filer exists."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the filer."},
			{Name: "tickers", Type: proto.ColumnType_JSON, Description: "Ticker of the filer."},
			{Name: "exchanges", Type: proto.ColumnType_JSON, Description: "Exchanges on which the filer trades."},
			{Name: "ein", Type: proto.ColumnType_STRING, Transform: transform.FromField("EIN"), Description: "EIN (Employer Identification Number) of the filer."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the filer."},
			{Name: "website", Type: proto.ColumnType_STRING, Description: "Website of the filer."},
			{Name: "investor_website", Type: proto.ColumnType_STRING, Description: "Investor website of the filer."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Category of the filer."},
			{Name: "fiscal_year_end", Type: proto.ColumnType_STRING, Description: "Fiscal year end of the filer."},
			{Name: "state_of_incorporation", Type: proto.ColumnType_STRING, Description: "State of incorporation of the filer."},
			{Name: "state_of_incorporation_description", Type: proto.ColumnType_STRING, Description: "State of incorporation description of the filer."},
			{Name: "addresses", Type: proto.ColumnType_STRING, Description: "Addresses of the filer."},
			{Name: "phone", Type: proto.ColumnType_STRING, Description: "Phone of the filer."},
			{Name: "flags", Type: proto.ColumnType_STRING, Description: "Flags of the filer."},
			{Name: "former_names", Type: proto.ColumnType_JSON, Description: "Former names of the filer."},
		},
	}
}

func listSecFiler(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO: move client init to main() in main.go
	logger := plugin.Logger(ctx)
	apiKey := os.Getenv("IEX_API_KEY")
	logger.Info("IEX API key = ", apiKey)
	if apiKey == "" {
		panic("No IEX API Key found")
	}
	client := edgar.NewClient(apiKey)
	quals := d.KeyColumnQuals
	cik := quals["cik"].GetStringValue()
	filer, err := client.GetSubmissions(cik)
	if err != nil {
		logger.Error("tableSecFilers.listSecFiler", "query_error", err)
		return nil, err
	}
	d.StreamListItem(ctx, filer)
	return nil, nil
}

// transformStringToInt
func transformStringToInt(ctx context.Context, td *transform.TransformData) (interface{}, error) {
	sicString, err := td.Value.(*string)
	if !err {
		panic(*(td.Value.(*string)))
	}
	if len(*sicString) == 0 {
		return nil, nil
	}
	return sicString, nil
}

// transformCIK
func transformCIK(ctx context.Context, td *transform.TransformData) (interface{}, error) {
	shortCik, err := td.Value.(*string)
	if !err {
		panic(*(td.Value.(*string)))
	}

	// sometimes EDGAR stores CIKs as 9 digit strings with the preceding zeros and other times it
	// does not. The behaviour is inconsistent.
	var cik string = *shortCik
	var cikLen int = len(*shortCik)
	const standardCikLen int = 10
	if cikLen == standardCikLen {
		return &cik, nil
	} else if cikLen < standardCikLen {
		var prefix string
		for i := 0; i < standardCikLen-cikLen; i++ {
			prefix += "0"
		}
		cik = prefix + cik
		return &cik, nil
	} else {
		panic("Invalid CIK: CIK longer than 10 characters.")
	}
}
