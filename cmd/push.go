/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/mysql"
	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/rtl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var (
	pushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push data into MySQL",
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
			if err := pushDataToMySQL(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("error")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVarP(&datafile, "datafile", "f", "data/output.gob", "gob data file")
	pushCmd.Flags().StringVarP(&dsn, "dsn", "d", "rtl:rtl@tcp(127.0.0.1:3306)/rtl", "MySQL DSN")
}

func pushDataToMySQL() error {
	db, err := mysql.New(mysql.SetDSN(dsn))
	if err != nil {
		return err
	}

	// Resolve the output file path
	fqpn, err := filepath.Abs(datafile)
	if err != nil {
		return err
	}
	f, err := os.Open(fqpn)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var data []*rtl.Record
	if err := dec.Decode(&data); err != nil {
		return err
	}

	fmt.Printf("Pushing %d records to MySQL\n", len(data))
	for i, r := range data {
		if err := db.Insert(&mysql.Record{Record: r}); err != nil {
			return err
		}
		if i%1000 == 0 {
			fmt.Printf(" ..%d.. ", i)
		}
	}
	fmt.Println(" ..done")

	return nil
}
