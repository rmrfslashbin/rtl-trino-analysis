/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/rmrfslashbin/rtl-trino-analysis/pkg/useragent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uapCmd represents the uap command
var uapCmd = &cobra.Command{
	Use:   "uap",
	Short: "user agent parser",
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
		if err := uap(); err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error")
		}
	},
}

func init() {
	rootCmd.AddCommand(uapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	uapCmd.Flags().String("tsv", "", "tsv file")
	viper.BindPFlag("tsv", uapCmd.Flags().Lookup("tsv"))
}

func uap() error {
	tsv := viper.GetString("tsv")
	if tsv == "" {
		return fmt.Errorf("tsv is required")
	}

	tsvFile, err := filepath.Abs(tsv)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tsvFile); os.IsNotExist(err) {
		return fmt.Errorf("tsv (%s) file does not exist", tsvFile)
	}

	fh, err := os.Open(tsvFile)
	if err != nil {
		return err
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		uaText, err := url.QueryUnescape(scanner.Text())
		if err != nil {
			return err
		}

		ua, err := useragent.Parse(uaText)
		if err != nil {
			return err
		}
		spew.Dump(ua)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nil
}
