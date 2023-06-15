package user

import (
	"fmt"
	"github.com/gabeduke/wio-cli-go/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUserCmd() *cobra.Command {
	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "Manage your wio user",
	}

	userCmd.AddCommand(newUserCreateCmd())
	userCmd.AddCommand(newConfigureCmd())

	return userCmd
}

func newUserCreateCmd() *cobra.Command {
	var userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("user")
			credentials := &credentials{}
			_, err := credentials.Create(logger)
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Println("Success!")
		},
	}

	userCreateCmd.PersistentFlags().StringP("Email", "e", "", "Email address")
	viper.BindPFlag("Email", userCreateCmd.PersistentFlags().Lookup("Email"))

	return userCreateCmd
}

func NewUserLoginCmd() *cobra.Command {
	var userLoginCmd = &cobra.Command{
		Use:     "login",
		Short:   "Login to Wio Server",
		Long:    "Login to the Wio Server and store the returned token in the configuration file.",
		Aliases: []string{"auth", "authenticate"},
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("user")
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

	userLoginCmd.PersistentFlags().StringP("Email", "e", "", "Email address")
	viper.BindPFlag("Email", userLoginCmd.PersistentFlags().Lookup("Email"))

	return userLoginCmd
}

func newConfigureCmd() *cobra.Command {
	// configureCmd represents the configure command
	var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Setup the Wio CLI Configuration file",
		Long: `The default configuration file is located at ~/.wio/config.json

The configuration file is used to store your 
* Wio API Access Token
* User Email address
* Server address
* Server IP

This command will Prompt you for the above information and store it in the configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("user")
			err := configure(logger)
			if err != nil {
				logger.Fatal(err)
			}
		},
	}
	return configureCmd
}
