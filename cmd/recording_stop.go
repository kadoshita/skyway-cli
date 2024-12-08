package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/kadoshita/skyway-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// recordingStopCmd represents the stop command
var recordingStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop recording by delete a recording session",
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

		response, err := internal.DeleteSession(channelId, sessionId, token, url)
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
	recordingCmd.AddCommand(recordingStopCmd)

	recordingStopCmd.Flags().String("channel-id", "", "Channel ID")
	recordingStopCmd.Flags().String("session-id", "", "RecordingSession ID")
}
