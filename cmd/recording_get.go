package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kadoshita/skyway-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// recordingGetCmd represents the get command
var recordingGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a recording session",
	Run: func(cmd *cobra.Command, args []string) {
		appId := viper.GetString("skyway.app_id")
		secretKey := viper.GetString("skyway.secret_key")
		url := viper.GetString("skyway.recording.url")

		channelId, err := cmd.Flags().GetString("channel-id")
		cobra.CheckErr(err)

		sessionId, err := cmd.Flags().GetString("session-id")
		cobra.CheckErr(err)

		pretty, err := cmd.Flags().GetBool("pretty")
		cobra.CheckErr(err)

		token, err := GenerateAdminToken(appId, secretKey, 3600, []string{})
		cobra.CheckErr(err)

		response, err := internal.GetSession(channelId, sessionId, token, url)
		cobra.CheckErr(err)

		if pretty {
			jsonString, err := json.MarshalIndent(response, "", "  ")
			cobra.CheckErr(err)

			fmt.Println(string(jsonString))
		} else {
			jsonString, err := json.Marshal(response)
			cobra.CheckErr(err)

			fmt.Println(string(jsonString))
		}
	},
}

func init() {
	recordingCmd.AddCommand(recordingGetCmd)

	recordingGetCmd.Flags().String("channel-id", "", "Channel ID")
	recordingGetCmd.Flags().String("session-id", "", "RecordingSession ID")
}
