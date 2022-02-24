/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/rtl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var (
	datafile string

	statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
			if err := getStats(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Fatal("error")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.PersistentFlags().StringVarP(&datafile, "datafile", "f", "data/output.gob", "gob data file")
}

func getStats() error {
	// Resolve the output file path
	fqpn, err := filepath.Abs(outfile)
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

	fmt.Printf("%d records in file %s\n", len(data), fqpn)

	ipAddrs := make(map[string]int)
	for _, r := range data {
		ipAddrs[string(r.ClientIP.IP.String())]++
	}
	for k, v := range ipAddrs {
		fmt.Printf("%s: %d\n", k, v)
	}

	return nil
}
