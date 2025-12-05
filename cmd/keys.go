package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/api"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
	"strconv"
)

var keysCmnd = &cobra.Command{
	Use:   "keys",
	Short: "For api key management",
}

var listKeys = &cobra.Command{
	Use:   "list",
	Short: "List api keys",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		resp, error := client.ListKeys()
		if error != nil {
			fmt.Printf("could not get api keys: %v", error)
			return
		}

		for _, k := range resp {
			fmt.Printf("id: %v , name: %v, key: %v, created: %v, expires: %v, last used: %v\n", k.ID, k.Name, k.Key[:8], k.CreatedAt, k.ExpiresAt, k.LastUsed)
		}

	},
}

var generateKey = &cobra.Command{
	Use:   "generate <name>",
	Short: "Create new api key",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		name := args[0]

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		resp, error := client.GenerateKey(name)
		if error != nil {
			fmt.Printf("could not generate api key: %v", error)
			return
		}

		fmt.Printf("api key %v generated, key value: %v , please save this key now (you won't be able to see it again). key will expire %v,", resp.Name, resp.Key, resp.ExpiresAt)

	},
}

var revokeKey = &cobra.Command{
	Use:   "revoke <id>",
	Short: "revoke api key by id",
	Args:  cobra.ExactArgs(1),
	Run: func(command *cobra.Command, args []string) {
		id := args[0]

		integer, err := strconv.Atoi(id)
		if err != nil {
			fmt.Printf("Error, invalid ID, please check you've supplied an integer as argument: %v\n", err)
			return
		}

		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")
		client := api.NewClient(baseUrl, key)

		error := client.RevokeKey(integer)
		if error != nil {
			fmt.Printf("could not delete api key, %v", error)
			return
		}

		fmt.Printf("api key of id %v , has been successfully removed", integer)

	},
}

func init() {
	rootCmd.AddCommand(keysCmnd)
	keysCmnd.AddCommand(listKeys)
	keysCmnd.AddCommand(generateKey)
	keysCmnd.AddCommand(revokeKey)
}
