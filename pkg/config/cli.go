package config

import (
	"github.com/gabeduke/wio-cli-go/pkg/util"
	"github.com/spf13/cobra"
)

func NewConfigureCmd(cfgFile string) *cobra.Command {
	// configureCmd represents the configure command
	var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Setup the Wio CLI Configuration file",
		Long: `The default configuration file is located at ~/.wio/config.json

The configuration file is used to store your 
* Wio API Access Token
* User email address
* Server address
* Server IP

This command will Prompt you for the above information and store it in the configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger := util.CreateNamedLogger("config")
			err := configure(logger, cfgFile)
			if err != nil {
				logger.Fatal(err)
			}
		},
	}
	return configureCmd
}
