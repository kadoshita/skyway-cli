package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kadoshita/skyway-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetChannel(channelId string, token string) (string, error) {
	return "", nil
}

// channelGetCmd represents the get command
var channelGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a channel",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.channel.url", cmd.Flags().Lookup("url"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		appId := viper.GetString("skyway.app_id")
		secretKey := viper.GetString("skyway.secret_key")
		url := viper.GetString("skyway.channel.url")

		pretty, err := cmd.Flags().GetBool("pretty")
		cobra.CheckErr(err)

		token, err := GenerateAdminToken(appId, secretKey, 3600, []string{})
		cobra.CheckErr(err)

		channel, err := internal.GetChannel(id, "", token, url)
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
	channelCmd.AddCommand(channelGetCmd)

	channelGetCmd.Flags().String("url", "https://channel.skyway.ntt.com/v1/json-rpc", "SkyWay Channel API URL. This option can also be set by the skyway.channel.url configuration or the SKYWAY_CHANNEL_URL environment variable.")
}
