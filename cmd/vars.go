package cmd

import (
	"github.com/sirupsen/logrus"
)

const PROGRAM_NAME = "rtl-trino-analysis"

var (
	//datafile string
	//trinoDSN string
	//mysqlDSN string
	//hostname string
	//outfile  string
	//geodb    string
	//geo      *geoip.Config
	log           *logrus.Logger
	cfgFile       string
	homeConfigDir string
	Version       string
)
