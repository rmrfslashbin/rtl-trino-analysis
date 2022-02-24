package cmd

import "github.com/rmrfslashbin/rtl-trino-analysis/pkg/geoip"

var (
	datafile string
	trinoDSN string
	mysqlDSN string
	hostname string
	outfile  string
	geodb    string
	geo      *geoip.Config
)
