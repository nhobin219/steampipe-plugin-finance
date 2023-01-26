package finance

import (
	"context"
	"os"

	"github.com/turbot/steampipe-plugin-finance/pkg/edgar"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableSecFilings(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "sec_filings",
		Description: "US public company filings from the SEC Edgar database.",
		List: &plugin.ListConfig{
			Hydrate:    listSecFilings,
			KeyColumns: plugin.SingleColumn("cik"),
		},
		Columns: []*plugin.Column{
			{Name: "cik", Type: proto.ColumnType_STRING, Transform: transform.FromField("CIK").Transform(transformCIK), Description: "CIK (Central Index Key) of the filer."},
			{Name: "accession_number", Type: proto.ColumnType_STRING, Description: "Accession number of the filing."},
			{Name: "filing_date", Type: proto.ColumnType_STRING, Description: "Filing date of the filing."},
			{Name: "report_date", Type: proto.ColumnType_STRING, Description: "Report date of the company."},
			{Name: "acceptance_date_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("AcceptanceDateTime"), Description: "Acceptance datetime of the filing."},
			{Name: "act", Type: proto.ColumnType_STRING, Description: "Act of the filing."},
			{Name: "form", Type: proto.ColumnType_STRING, Description: "Form of the filing."},
			{Name: "file_number", Type: proto.ColumnType_STRING, Description: "File number of the filing."},
			{Name: "film_number", Type: proto.ColumnType_STRING, Description: "Film number of the filing."},
			{Name: "items", Type: proto.ColumnType_STRING, Description: "Items of the filing."},
			{Name: "size", Type: proto.ColumnType_STRING, Description: "Size of the filing."},
			{Name: "is_xbrl", Type: proto.ColumnType_INT, Transform: transform.FromField("IsXBRL"), Description: "Whether or not the filing is in XBRL format."},
			{Name: "is_inline_xbrl", Type: proto.ColumnType_INT, Transform: transform.FromField("IsInlineXBRL"), Description: "Whether or not the filing is in inline XBRL format."},
			{Name: "primary_document", Type: proto.ColumnType_STRING, Description: "Primary document of the filing."},
			{Name: "primary_doc_description", Type: proto.ColumnType_STRING, Description: "Primary document description."},
		},
	}
}

func listSecFilings(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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
		logger.Error("tableSecFilings.listSecFilings", "query_error", err)
		return nil, err
	}
	for idx := range *filer.Filings.Recent.AccessionNumber {
		filing := edgar.Filing{}
		filing.CIK = filer.CIK
		// for idx, _ := range *filer.Filings.Recent.AccessionNumber {
		// 	logger.Error((*filer.Filings.Recent.AccessionNumber)[idx])
		// }
		if filer.Filings.Recent.AccessionNumber != nil {
			filing.AccessionNumber = &(*filer.Filings.Recent.AccessionNumber)[idx]
		}
		if filer.Filings.Recent.FilingDate != nil {
			filing.FilingDate = &(*filer.Filings.Recent.FilingDate)[idx]
		}
		if filer.Filings.Recent.ReportDate != nil {
			filing.ReportDate = &(*filer.Filings.Recent.ReportDate)[idx]
		}
		if filer.Filings.Recent.AcceptanceDateTime != nil {
			filing.AcceptanceDateTime = &(*filer.Filings.Recent.AcceptanceDateTime)[idx]
		}
		if filer.Filings.Recent.Act != nil {
			filing.Act = &(*filer.Filings.Recent.Act)[idx]
		}
		if filer.Filings.Recent.Form != nil {
			filing.Form = &(*filer.Filings.Recent.Form)[idx]
		}
		if filer.Filings.Recent.FileNumber != nil {
			filing.FileNumber = &(*filer.Filings.Recent.FileNumber)[idx]
		}
		if filer.Filings.Recent.FilmNumber != nil {
			filing.FilmNumber = &(*filer.Filings.Recent.FilmNumber)[idx]
		}
		if filer.Filings.Recent.Items != nil {
			filing.Items = &(*filer.Filings.Recent.Items)[idx]
		}
		if filer.Filings.Recent.Size != nil {
			filing.Size = &(*filer.Filings.Recent.Size)[idx]
		}
		if filer.Filings.Recent.IsXBRL != nil {
			filing.IsXBRL = &(*filer.Filings.Recent.IsXBRL)[idx]
		}
		if filer.Filings.Recent.IsInlineXBRL != nil {
			filing.IsInlineXBRL = &(*filer.Filings.Recent.IsInlineXBRL)[idx]
		}
		if filer.Filings.Recent.PrimaryDocument != nil {
			filing.PrimaryDocument = &(*filer.Filings.Recent.PrimaryDocument)[idx]
		}
		if filer.Filings.Recent.PrimaryDocDescription != nil {
			filing.PrimaryDocDescription = &(*filer.Filings.Recent.PrimaryDocDescription)[idx]
		}
		d.StreamListItem(ctx, &filing)
	}
	return nil, nil
}
