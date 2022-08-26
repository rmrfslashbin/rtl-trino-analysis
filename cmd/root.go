/*
Copyright © 2022 rmrfslashbin@sigler.io

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
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   PROGRAM_NAME,
	Short: "Realtime traffic analysis",
	Long: `rtl-trino-analysis is a realtime traffic analysis tool.
	
	The tool helps pull data from AWS Glue via Trino and analyze it.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	var err error

	// Find home directory.
	homeConfigDir, err = os.UserConfigDir()
	cobra.CheckErr(err)
	homeConfigDir = path.Join(homeConfigDir, PROGRAM_NAME)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s/config.yaml)", homeConfigDir))
	rootCmd.PersistentFlags().String("loglevel", "info", "log level (debug, info, warn, error, fatal, panic)")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(homeConfigDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	viper.ReadInConfig()

	// Set log level
	switch viper.GetString("loglevel") {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Log level set to debug")
	case "info":
		log.SetLevel(logrus.InfoLevel)
		log.Info("Log level set to info")
	case "warn":
		log.SetLevel(logrus.WarnLevel)
		log.Warn("Log level set to warn")
	case "error":
		log.SetLevel(logrus.ErrorLevel)
		log.Error("Log level set to error")
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
		log.Fatal("Log level set to fatal")
	case "panic":
		log.SetLevel(logrus.PanicLevel)
		log.Panic("Log level set to panic")
	default:
		log.SetLevel(logrus.InfoLevel)
		log.Info("Log level set to info")
	}
}
