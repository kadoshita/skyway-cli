package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// recordingCmd represents the recording command
var recordingCmd = &cobra.Command{
	Use:   "recording",
	Short: "Audio and video recording",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.app_id", cmd.PersistentFlags().Lookup("app-id"))
		viper.BindPFlag("skyway.secret_key", cmd.PersistentFlags().Lookup("secret-key"))
		viper.BindPFlag("skyway.recording.url", cmd.PersistentFlags().Lookup("url"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("recording called")
	},
}

func init() {
	rootCmd.AddCommand(recordingCmd)

	recordingCmd.PersistentFlags().String("app-id", "", "SkyWay App ID. This option can also be set by the skyway.app_id configuration or the SKYWAY_APP_ID environment variable.")
	recordingCmd.PersistentFlags().String("secret-key", "", "SkyWay Secret Key. This option can also be set by the skyway.secret_key configuration or the SKYWAY_SECRET_KEY environment variable.")
	recordingCmd.PersistentFlags().BoolP("pretty", "p", false, "Pretty print JSON")
	recordingCmd.PersistentFlags().String("url", "", "SkyWay Recording API URL. This option can also be set by the skyway.recording.url configuration or the SKYWAY_RECORDING_URL environment variable.")
}
