package user

import (
	"fmt"
	"github.com/gabeduke/wio-cli-go/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUserCmd() *cobra.Command {
	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "Manage user settings",
	}

	userCmd.AddCommand(newUserCreateCmd())
	userCmd.AddCommand(newUserLoginCmd())

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
			err := resp.login(logger)
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

