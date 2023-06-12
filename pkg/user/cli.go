package user

import (
	"fmt"
	"github.com/gabeduke/wio-cli-go/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUserCmd(cfg string) *cobra.Command {
	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "Manage user settings",
	}

	userCmd.AddCommand(newUserCreateCmd())
	userCmd.AddCommand(newUserLoginCmd())
	userCmd.AddCommand(newConfigureCmd(cfg))

	return userCmd
}

func newUserCreateCmd() *cobra.Command {
	var userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			logger := util.CreateNamedLogger("user")
			credentials := &credentials{}
			_, err := credentials.create(logger)
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Println("Success!")
		},
	}

	userCreateCmd.PersistentFlags().StringP("email", "e", "", "email address")
	viper.BindPFlag("email", userCreateCmd.PersistentFlags().Lookup("email"))

	return userCreateCmd
}

func newUserLoginCmd() *cobra.Command {
	var userLoginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to Wio Server",
		Long:  "Login to the Wio Server and store the returned token in the configuration file.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := util.CreateNamedLogger("user")
			resp := &LoginResponse{}
			err := resp.Login(logger)
			if err != nil {
				logger.Fatal(err)
			}

			viper.Set("token", resp.Token)

			err = viper.WriteConfig()
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Printf("Login successful. Writing token to %s\n", viper.ConfigFileUsed())
		},
	}

	userLoginCmd.PersistentFlags().StringP("email", "e", "", "email address")
	viper.BindPFlag("email", userLoginCmd.PersistentFlags().Lookup("email"))

	return userLoginCmd
}

func newConfigureCmd(cfgFile string) *cobra.Command {
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
			logger := util.CreateNamedLogger("user")
			err := configure(logger, cfgFile)
			if err != nil {
				logger.Fatal(err)
			}
		},
	}
	return configureCmd
}
