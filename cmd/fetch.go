/*
Copyright © 2022 Robert Sigler <sigler@improvisedscience.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/fetch"
	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/geoip"
	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/rtl"
	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/useragent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetches data from Trino and writes it to a gob file for furthe processing",
	Run: func(cmd *cobra.Command, args []string) {
		// Catch errors
		var err error
		defer func() {
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("main crashed")
			}
		}()
		if err := fetchData(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.PersistentFlags().StringP("trinodsn", "d", "", "Trino DNS")
	viper.BindPFlag("trinodsn", fetchCmd.Flags().Lookup("dsn"))

	fetchCmd.PersistentFlags().StringP("geoipdb", "g", "", "Path to GeoIP database")
	viper.BindPFlag("geoipdb", fetchCmd.Flags().Lookup("geoip"))

	fetchCmd.PersistentFlags().StringP("datafile", "o", "", "GOB Output file")
	viper.BindPFlag("datafile", fetchCmd.Flags().Lookup("datafile"))

	fetchCmd.PersistentFlags().StringP("hostname", "n", "", "Hostname to query")
	viper.BindPFlag("hostname", fetchCmd.Flags().Lookup("hostname"))
}

func fetchData() error {
	trinodns := viper.GetString("trinodsn")
	if trinodns == "" {
		return fmt.Errorf("trinodsn is required")
	}

	geoipdb := viper.GetString("geoipdb")
	if geoipdb == "" {
		return fmt.Errorf("geoipdb is required")
	}

	datafile := viper.GetString("datafile")
	if datafile == "" {
		return fmt.Errorf("datafile is required")
	}

	hostname := viper.GetString("hostname")
	if hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	log.WithFields(logrus.Fields{
		"trinodsn": trinodns,
		"geoipdb":  geoipdb,
		"datafile": datafile,
		"hostname": hostname,
	}).Debug("fetching data")

	geo, err := geoip.New(geoip.SetGeoDB(geoipdb))
	if err != nil {
		return err
	}

	// Get a new Trino database connection
	trino, err := fetch.New(fetch.SetDSN(trinodns))
	if err != nil {
		return err
	}
	defer trino.Close()

	log.Debug("Connected to Trino")

	// Query to execute
	query := "SELECT * FROM hive.cfrtl.rtl WHERE year='2022' AND month='7'"
	// "SELECT * FROM hive.cfrtl.rtl WHERE year='2022' AND month='2' AND HOST LIKE ?"

	log.WithFields(logrus.Fields{
		"query": query,
	}).Debug("Executing query")
	// Execute the query
	rows, err := trino.DB.Queryx(query)
	//rows, err := trino.DB.Queryx(query, fmt.Sprintf("%%%s", hostname))
	if err != nil {
		return err
	}
	defer rows.Close()
	log.Debug("Query executed. Processing results...")

	// Slice to hold the results
	var records []rtl.Record

	// Iterate over the results
	count := 0
	for rows.Next() {
		// fetch.Entry is a struct that is compatible with the data returned from Trino
		var entry fetch.Entry

		// Get a row
		if err := rows.StructScan(&entry); err != nil {
			return err
		}

		record, err := processRecord(geo, &entry)
		if err != nil {
			return err
		}

		// Build a properly typed reccord
		records = append(records, *record)
		count++
		if count%1000 == 0 {
			log.WithFields(logrus.Fields{
				"count": count,
			}).Debug("Processed record")
		}
	}
	log.WithFields(logrus.Fields{
		"count": count,
	}).Debug("Processed all records")

	// Write the records to a file
	if fqpn, err := writeData(datafile, &records); err != nil {
		return err
	} else {
		// Print the number of records and output filename
		log.WithFields(logrus.Fields{
			"count": len(records),
			"file":  fqpn,
		}).Info("Wrote data")
	}

	geo.Close()

	return nil
}

func processRecord(geo *geoip.Config, entry *fetch.Entry) (*rtl.Record, error) {
	// Convert the timestamp string to a time.Time
	epoch, err := strconv.ParseInt(strings.Replace(entry.Timestamp, ".", "", 1), 10, 64)
	if err != nil {
		return nil, err
	}
	timestamp := time.UnixMilli(epoch)

	// Convert the year, month, and day strings to ints
	year, err := strconv.Atoi(entry.Year)
	if err != nil {
		return nil, err
	}
	month, err := strconv.Atoi(entry.Month)
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(entry.Day)
	if err != nil {
		return nil, err
	}

	geodata, err := geo.Lookup(net.ParseIP(entry.ClientIP))
	if err != nil {
		return nil, err
	}

	uaparser, err := useragent.Parse(entry.UserAgent)
	if err != nil {
		return nil, err
	}

	return &rtl.Record{
		Timestamp:                timestamp,
		ClientIPAddr:             entry.ClientIP,
		Status:                   entry.Status,
		Bytes:                    entry.Bytes,
		Method:                   entry.Method,
		Protocol:                 entry.Protocol,
		Host:                     entry.Host,
		UriStem:                  entry.UriStem,
		EdgeLocation:             entry.EdgeLocation,
		EdgeRequestID:            entry.EdgeRequestID,
		HostHeader:               entry.HostHeader,
		TimeTaken:                entry.TimeTaken,
		ProtoVersion:             entry.ProtoVersion,
		IPVersion:                entry.IPVersion,
		Referer:                  entry.Referer,
		Cookie:                   entry.Cookie,
		UriQuery:                 entry.UriQuery,
		EdgeResponseResultType:   entry.EdgeResponseResultType,
		SslProtocol:              entry.SslProtocol,
		SslCipher:                entry.SslCipher,
		EdgeResultType:           entry.EdgeResultType,
		ContentType:              entry.ContentType,
		ContentLength:            entry.ContentLength,
		EdgeDetailedResultType:   entry.EdgeDetailedResultType,
		Country:                  entry.Country,
		CacheBehaviorPathPattern: entry.CacheBehaviorPathPattern,
		Year:                     year,
		Month:                    month,
		Day:                      day,
		ClientIP:                 geodata,
		UserAgent:                uaparser,
	}, nil
}

func writeData(datafile string, records *[]rtl.Record) (*string, error) {
	// Resolve the output file path
	fqpn, err := filepath.Abs(datafile)
	if err != nil {
		return nil, err
	}

	// Create the output file
	f, err := os.Create(fqpn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Encode the records to the file
	enc := gob.NewEncoder(f)
	if err := enc.Encode(records); err != nil {
		return nil, err
	}

	return &fqpn, nil
}
