package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
	"reflect"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(command *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		err := config.Set(key, value)
		if err != nil {
			fmt.Printf("could not set: %v\n", err)
			return
		}

		fmt.Printf("New %v has been set", key)

	},
}

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		key := args[0]

		value, err := config.Get(key)
		if err != nil {
			fmt.Printf("could not get: %v\n", err)
			return
		}

		fmt.Printf("Current %v is: %v\n", key, value)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all current config values",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		currentConfig, err := config.Load()
		if err != nil {
			fmt.Printf("could not load config file: %v\n", err)
			return
		}

		values := reflect.ValueOf(*currentConfig)
		types := values.Type()

		for i := 0; i < values.NumField(); i++ {
			fmt.Println(types.Field(i).Tag.Get("json"), ": ", values.Field(i))
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(getCmd)
	configCmd.AddCommand(listCmd)
}
