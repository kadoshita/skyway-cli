package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// channelCmd represents the channel command
var channelCmd = &cobra.Command{
	Use:   "channel",
	Short: "Channel operations",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.app_id", cmd.PersistentFlags().Lookup("app-id"))
		viper.BindPFlag("skyway.secret_key", cmd.PersistentFlags().Lookup("secret-key"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("channel called")
	},
}

func init() {
	rootCmd.AddCommand(channelCmd)

	channelCmd.PersistentFlags().String("app-id", "", "SkyWay App ID. This option can also be set by the skyway.app_id configuration or the SKYWAY_APP_ID environment variable.")
	channelCmd.PersistentFlags().String("secret-key", "", "SkyWay Secret Key. This option can also be set by the skyway.secret_key configuration or the SKYWAY_SECRET_KEY environment variable.")
	channelCmd.PersistentFlags().BoolP("pretty", "p", false, "Pretty print JSON")
}
