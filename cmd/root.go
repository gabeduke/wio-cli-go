/*
Package cmd

Copyright © 2023 Gabriel Duke gabeduke@gmail.com

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
	"errors"
	"fmt"
	"github.com/gabeduke/wio-cli-go/pkg/nodes"
	"github.com/gabeduke/wio-cli-go/pkg/user"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel logLevelEnum
)

type logLevelEnum string

const (
	logEnumInfoLevel  logLevelEnum = "info"
	logEnumDebugLevel logLevelEnum = "debug"
	logEnumWarnLevel  logLevelEnum = "warn"
	logEnumErrorLevel logLevelEnum = "error"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *logLevelEnum) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *logLevelEnum) Set(v string) error {
	switch v {
	case "info", "debug", "warn", "error":
		*e = logLevelEnum(v)
		return nil
	default:
		return errors.New(`must be one of "info", "debug", "warn", "error"`)
	}
}

// Type is only used in help text
func (e *logLevelEnum) Type() string {
	return "logLevelEnum"
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "wio",
	Aliases: []string{"wio-cli-go", "wio-cli"},
	Short:   "A CLI for interacting with the Wio Link IoT Platform",
	Long: `This is a rewrite of the Wio CLI in Go. It is a CLI for interacting with the Wio Link IoT Platform.

The original Wio CLI was written in Python and can be found here: https://github.com/Seeed-Studio/wio-cli

This CLI is a work in progress and is not yet feature complete. It is being developed by Gabriel Duke (gabeduke@gmail.com)
The goal of this rewrite is to make the CLI more performant and easier to maintain.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wio.json)")
	rootCmd.PersistentFlags().VarP(&logLevel, "log-level", "l", `log level: "info", "debug", "warn", "error" (default is warn)`)
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(user.NewUserCmd(cfgFile))
	rootCmd.AddCommand(nodes.NewNodesCmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var logDebugMessages []string
	var logFatalMessages []string
	var logErrorMessages []string

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			logFatalMessages = append(logFatalMessages, fmt.Sprintf("Error finding home directory: %s", err))
		}

		// Search config in home directory with name ".wio" (without extension).
		viper.AddConfigPath(home + "/.wio")
		viper.SetConfigType("json")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		logErrorMessages = append(logErrorMessages, fmt.Sprintf("Error reading config file: %s", err))
	} else {
		logDebugMessages = append(logDebugMessages, fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
	}

	initLogger()

	// Log messages
	for _, message := range logDebugMessages {
		log.Infof(message)
	}

	for _, message := range logFatalMessages {
		log.Fatalf(message)
	}
}

func initLogger() {
	logLevel := viper.GetString("logLevel")
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}
