package finance

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-finance/pkg/edgar"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableCompanies(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "companies",
		Description: "US public companies from the SEC Edgar database.",
		List: &plugin.ListConfig{
			Hydrate: listCompanies,
		},
		Columns: []*plugin.Column{
			{Name: "symbol", Type: proto.ColumnType_STRING, Description: "Symbol of the company."}, // necessary
			{Name: "exchange", Type: proto.ColumnType_STRING, Description: "Exchange of the company."},
			{Name: "exchange_suffix", Type: proto.ColumnType_STRING, Transform: transform.FromField("ExchangeSuffix	"), Description: "Exchange suffix of the company."},
			{Name: "exchange_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("ExchangeName"), Description: "Exchange name of the company."},
			{Name: "exchange_segment", Type: proto.ColumnType_STRING, Description: "Exchange segment of the company."},
			{Name: "exchange_segment_name", Type: proto.ColumnType_STRING, Description: "Exchange segment name of the company."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the company."},
			{Name: "date", Type: proto.ColumnType_STRING, Description: "Date of the company."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the company."},
			{Name: "iex_id", Type: proto.ColumnType_STRING, Description: "IexId of the company."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "Region of the company."},
			{Name: "currency", Type: proto.ColumnType_STRING, Description: "Currency of the company."},
			{Name: "is_enabled", Type: proto.ColumnType_BOOL, Description: "Whether or not the company is enabled."},
			{Name: "figi", Type: proto.ColumnType_STRING, Transform: transform.FromField("FIGI"), Description: "Financial Instrument Global Identifier of the company."},
			{Name: "cik", Type: proto.ColumnType_STRING, Transform: transform.FromField("CIK"), Description: "Central Index Key of the company."}, // necessary
			{Name: "lei", Type: proto.ColumnType_STRING, Transform: transform.FromField("LEI"), Description: "Legal Entity Identifier of the company."},
		},
	}
}

func listCompanies(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	apiKey := os.Getenv("IEX_API_KEY")
	logger.Info("IEX API key = ", apiKey)
	if apiKey == "" {
		panic("No IEX API Key found")
	}
	client := edgar.NewClient(apiKey)
	companies, err := client.GetPublicCompanies()
	if err != nil {
		logger.Error("companies.listCompanies", "query_error", err)
		return nil, err
	}
	for _, c := range *companies {
		d.StreamListItem(ctx, c)
	}
	return nil, nil
}
