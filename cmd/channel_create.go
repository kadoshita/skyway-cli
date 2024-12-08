package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kadoshita/skyway-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// channelCreateCmd represents the create command
var channelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a channel",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.channel.url", cmd.Flags().Lookup("url"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		channelName, err := cmd.Flags().GetString("name")
		cobra.CheckErr(err)
		metadata, err := cmd.Flags().GetString("metadata")
		cobra.CheckErr(err)

		appId := viper.GetString("skyway.app_id")
		secretKey := viper.GetString("skyway.secret_key")
		url := viper.GetString("skyway.channel.url")

		pretty, err := cmd.Flags().GetBool("pretty")
		cobra.CheckErr(err)

		token, err := GenerateAdminToken(appId, secretKey, 3600, []string{})
		cobra.CheckErr(err)

		channel, err := internal.CreateChannel(channelName, metadata, token, url)
		cobra.CheckErr(err)

		if pretty {
			jsonString, err := json.MarshalIndent(channel, "", "  ")
			cobra.CheckErr(err)

			fmt.Println(string(jsonString))
		} else {
			jsonString, err := json.Marshal(channel)
			cobra.CheckErr(err)

			fmt.Println(string(jsonString))
		}
	},
}

func init() {
	channelCmd.AddCommand(channelCreateCmd)

	channelCreateCmd.Flags().String("metadata", "", "Channel metadata")
	channelCreateCmd.Flags().String("name", "", "Channel name")
	channelCreateCmd.Flags().String("url", "https://channel.skyway.ntt.com/v1/json-rpc", "SkyWay Channel API URL. This option can also be set by the skyway.channel.url configuration or the SKYWAY_CHANNEL_URL environment variable.")
}
