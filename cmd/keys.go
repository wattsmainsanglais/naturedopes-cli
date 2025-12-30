package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/api"
	"github.com/wattsmainsanglais/naturedopes-cli/pkg/config"
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Key", "Created", "Expires", "Last Used", "Revoked"})
		for _, k := range resp {
			lastUsedStr := "Never"
			if k.LastUsed != nil {
				lastUsedStr = *k.LastUsed
			}
			table.Append([]string{fmt.Sprintf("%d", k.ID), k.Name, fmt.Sprintf("%v...", k.Key[:8]), k.CreatedAt, k.ExpiresAt, lastUsedStr, fmt.Sprintf("%v", k.Revoked)})
		}
		table.Render()

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

		fmt.Printf("api key %v generated, key value: %v , please save this key now (you won't be able to see it again). key will expire %v,\n", resp.Name, resp.Key, resp.ExpiresAt)

	},
}

var revokeKey = &cobra.Command{
	Use:   "revoke",
	Short: "revoke the configured api key",
	Args:  cobra.ExactArgs(0),
	Run: func(command *cobra.Command, args []string) {
		baseUrl, _ := config.Get("api-url")
		key, _ := config.Get("api-key")

		if key == "" {
			fmt.Println("Error: No API key configured. Use 'config set api-key <key>' first.")
			return
		}

		fmt.Printf("Are you sure you want to revoke your API key? This cannot be undone. (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "yes" {
			fmt.Println("Revoke cancelled")
			return
		}

		client := api.NewClient(baseUrl, key)

		error := client.RevokeKey()
		if error != nil {
			fmt.Printf("could not delete api key: %v\n", error)
			return
		}

		fmt.Println("Your API key has been successfully revoked")
	},
}

func init() {
	rootCmd.AddCommand(keysCmnd)
	keysCmnd.AddCommand(listKeys)
	keysCmnd.AddCommand(generateKey)
	keysCmnd.AddCommand(revokeKey)
}
