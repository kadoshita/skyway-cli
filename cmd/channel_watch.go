package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/kadoshita/skyway-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const tokenTempl = `{
  "jti": "JTI_PLACEHOLDER",
  "iat": 0,
  "exp": 0,
  "version": 2,
  "scope": {
    "app": {
      "id": "APP_ID_PLACEHOLDER",
      "actions": [
        "read"
      ],
      "turn": true,
      "channels": [
        {
          "id": "%s",
          "name": "%s",
          "actions": [
            "read"
          ],
          "members": []
        }
      ]
    }
  }
}`

// channelWatchCmd represents the watch command
var channelWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch channel events",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.rtc_api.url", cmd.Flags().Lookup("url"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		appId := viper.GetString("skyway.app_id")
		secretKey := viper.GetString("skyway.secret_key")
		url := viper.GetString("skyway.rtc_api.url")

		id, err := cmd.Flags().GetString("id")
		cobra.CheckErr(err)

		name, err := cmd.Flags().GetString("name")
		cobra.CheckErr(err)

		pretty, err := cmd.Flags().GetBool("pretty")
		cobra.CheckErr(err)

		token, err := GenerateToken(fmt.Sprintf(tokenTempl, id, name), appId, secretKey, 3*24*60*60, []string{})
		cobra.CheckErr(err)

		handleEvents := make(chan string)
		go func() {
			for {
				event := <-handleEvents

				if pretty {
					var buffer bytes.Buffer
					err := json.Indent(&buffer, []byte(event), "", "  ")
					cobra.CheckErr(err)

					fmt.Println(buffer.String())
				} else {
					fmt.Println(event)
				}
			}
		}()
		err = internal.SubscribeEvents(id, name, token, appId, url, handleEvents)
		cobra.CheckErr(err)
	},
}

func init() {
	channelCmd.AddCommand(channelWatchCmd)

	channelWatchCmd.Flags().String("id", "", "Channel id")
	channelWatchCmd.Flags().String("name", "", "Channel name")

	channelWatchCmd.Flags().String("url", "wss://rtc-api.skyway.ntt.com/ws", "SkyWay RTC API URL. This option can also be set by the skyway.rtc_api.url configuration or the SKYWAY_RTC_API_URL environment variable.")
}
