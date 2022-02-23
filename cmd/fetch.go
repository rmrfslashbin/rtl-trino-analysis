/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var (
	hostname string
	outfile  string
	dsn      string

	fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetches data from Trino and writes it to a gob file for furthe processing",
		Run: func(cmd *cobra.Command, args []string) {
			// Catch errors
			var err error
			defer func() {
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"error": err,
					}).Fatal("main crashed")
				}
			}()
			if err := runFetch(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("error")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(fetchCmd)

	fetchCmd.PersistentFlags().StringVarP(&hostname, "hostname", "n", "", "hostname")
	fetchCmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", "output.gob", "gob data file")
	fetchCmd.PersistentFlags().StringVarP(&dsn, "dsn", "d", "http://user@localhost:9080?catalog=hive&schema=cfrtl", "Trino DSN")
	fetchCmd.MarkPersistentFlagRequired("hostname")
}

func runFetch() error {
	// Get a new Trino database connection
	trino, err := fetch.New(fetch.SetDSN(dsn))
	if err != nil {
		return err
	}
	defer trino.Close()

	// Query to execute
	query := "SELECT * FROM hive.cfrtl.rtl WHERE year='2022' AND month='2' AND HOST LIKE ? LIMIT 10"

	// Execute the query
	rows, err := trino.DB.Queryx(query, fmt.Sprintf("%%%s", hostname))
	if err != nil {
		return err
	}
	defer rows.Close()

	// Slice to hold the results
	var records []fetch.Record

	// Iterate over the results
	for rows.Next() {
		// fetch.Entry is a struct that is compatible with the data returned from Trino
		var entry fetch.Entry

		// Get a row
		if err := rows.StructScan(&entry); err != nil {
			return err
		}

		// Convert the timestamp string to a time.Time
		epoch, err := strconv.ParseInt(strings.Replace(entry.Timestamp, ".", "", 1), 10, 64)
		if err != nil {
			return err
		}
		timestamp := time.UnixMilli(epoch)

		// Convert the year, month, and day strings to ints
		year, err := strconv.Atoi(entry.Year)
		if err != nil {
			return err
		}
		month, err := strconv.Atoi(entry.Month)
		if err != nil {
			return err
		}
		day, err := strconv.Atoi(entry.Day)
		if err != nil {
			return err
		}

		// Build a properly typed reccord
		records = append(records, fetch.Record{
			Timestamp:                timestamp,
			ClientIP:                 net.ParseIP(entry.ClientIP),
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
			UserAgent:                entry.UserAgent,
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
		})
	}

	// Resolve the output file path
	fqpn, err := filepath.Abs(outfile)
	if err != nil {
		return err
	}

	// Create the output file
	f, err := os.Create(fqpn)
	if err != nil {
		return err
	}
	defer f.Close()

	// Encode the records to the file
	enc := gob.NewEncoder(f)
	if err := enc.Encode(records); err != nil {
		return err
	}

	// Print the number of records and output filename
	fmt.Printf("wrote %d records to %s\n", len(records), fqpn)

	return nil
}
