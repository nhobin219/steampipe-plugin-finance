package finance

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-finance",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo(),
		DefaultConcurrency: &plugin.DefaultConcurrencyConfig{
			TotalMaxConcurrency: 10,
		},
		TableMap: map[string]*plugin.Table{
			"companies":    tableCompanies(ctx),
			"sec_filers":   tableSecFilers(ctx),
			"sec_filings":  tableSecFilings(ctx),
			"quote":        tableFinanceQuote(ctx),
			"quote_daily":  tableFinanceQuoteDaily(ctx),
			"quote_hourly": tableFinanceQuoteHourly(ctx),
		},
	}
	return p
}
