package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gabeduke/wio-cli-go/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type boardEnum string

var nodeName string
var boardType boardEnum

const (
	boardEnumNode boardEnum = "node"
	boardEnumLink boardEnum = "link"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *boardEnum) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *boardEnum) Set(v string) error {
	switch v {
	case "node", "link":
		*e = boardEnum(v)
		return nil
	default:
		return errors.New(`must be one of "link", "node"`)
	}
}

// Type is only used in help text
func (e *boardEnum) Type() string {
	return "boardEnum"
}

func NewNodesCmd() *cobra.Command {
	var nodesCmd = &cobra.Command{
		Use:   "nodes",
		Short: "Manage your wio nodes",
	}

	nodesCmd.AddCommand(newNodesRegisterCmd())
	nodesCmd.AddCommand(newNodesCreateCmd())
	nodesCmd.AddCommand(newNodesDeleteCmd())
	nodesCmd.AddCommand(newNodesListCmd())

	return nodesCmd
}

func newNodesRegisterCmd() *cobra.Command {
	var sn, key string
	var nodesRegisterCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a node",
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("nodes")
			err := RegisterNode()
			if err != nil {
				logger.Fatal(err)
			}
		},
	}

	nodesRegisterCmd.Flags().BoolP("create", "c", false, "Create a new node")
	nodesRegisterCmd.Flags().StringVarP(&sn, "sn", "s", "", "Serial number of the node")
	nodesRegisterCmd.Flags().StringVarP(&key, "key", "k", "", "Key of the node")
	nodesRegisterCmd.Flags().StringVarP(&nodeName, "name", "n", "", "Name of the node")
	nodesRegisterCmd.Flags().Var(&boardType, "board", `Wio Board type. allowed: "node", "link""`)
	viper.BindPFlag("create", nodesRegisterCmd.Flags().Lookup("create"))
	viper.BindPFlag("name", nodesRegisterCmd.Flags().Lookup("name"))
	viper.BindPFlag("board", nodesRegisterCmd.Flags().Lookup("board"))
	viper.BindPFlag("sn", nodesRegisterCmd.Flags().Lookup("sn"))
	viper.BindPFlag("key", nodesRegisterCmd.Flags().Lookup("key"))

	return nodesRegisterCmd
}

func newNodesCreateCmd() *cobra.Command {
	var nodesCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new node",
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("nodes")
			resp, err := CreateNode(nodeName, boardType)
			if err != nil {
				logger.Fatal(err)
			}

			logger.Infof("Node created: %s", resp.String())

			data, err := json.Marshal(resp)
			if err != nil {
				logger.Fatal(err)
			}
			fmt.Printf("%s\n", data)
		},
	}

	nodesCreateCmd.Flags().StringVarP(&nodeName, "name", "n", "", "Name of the node")
	nodesCreateCmd.Flags().Var(&boardType, "board", `Wio Board type. allowed: "node", "link""`)
	viper.BindPFlag("name", nodesCreateCmd.Flags().Lookup("name"))
	viper.BindPFlag("board", nodesCreateCmd.Flags().Lookup("board"))

	cobra.MarkFlagRequired(nodesCreateCmd.Flags(), "name")
	cobra.MarkFlagRequired(nodesCreateCmd.Flags(), "board")

	return nodesCreateCmd
}

func newNodesDeleteCmd() *cobra.Command {
	var sn string
	var nodesDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a node",
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("nodes")
			err := DeleteNode(sn)
			if err != nil {
				logger.Fatal(err)
			}

			logger.WithField("sn", sn).Info("Successfully deleted node")
			fmt.Println("Successfully deleted node: " + sn)
		},
	}

	nodesDeleteCmd.Flags().StringVarP(&sn, "sn", "s", "", "Serial number of the node")
	viper.BindPFlag("sn", nodesDeleteCmd.Flags().Lookup("sn"))

	cobra.MarkFlagRequired(nodesDeleteCmd.Flags(), "sn")

	return nodesDeleteCmd
}

func newNodesListCmd() *cobra.Command {
	var nodesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all of your nodes",
		Run: func(cmd *cobra.Command, args []string) {
			logger := internal.CreateNamedLogger("nodes")
			nodes, err := ListNodes()
			if err != nil {
				logger.Fatal(err)
			}

			fmt.Println(nodes)
		},
	}

	return nodesListCmd
}
