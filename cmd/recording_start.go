package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kadoshita/skyway-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var optionToService map[string]string = map[string]string{
	"gcs":    "GOOGLE_CLOUD_STORAGE",
	"aws":    "AMAZON_S3",
	"wasabi": "WASABI",
}

// recordingStartCmd represents the start command
var recordingStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start recording by create a recording session",
	Long: `Start recording by create a recording session.
Note: You can specify the single publication ID and content type.
This command cannot specify multiple publication IDs yet.`,
	Run: func(cmd *cobra.Command, args []string) {
		appId := viper.GetString("skyway.app_id")
		secretKey := viper.GetString("skyway.secret_key")
		url := viper.GetString("skyway.recording.url")

		channelId, err := cmd.Flags().GetString("channel-id")
		cobra.CheckErr(err)

		pretty, err := cmd.Flags().GetBool("pretty")
		cobra.CheckErr(err)

		publicationId, err := cmd.Flags().GetString("publication-id")
		cobra.CheckErr(err)

		contentType, err := cmd.Flags().GetString("content-type")
		cobra.CheckErr(err)

		outputServiceName, err := cmd.Flags().GetString("output-service")
		cobra.CheckErr(err)

		outputServiceConfigKey := "skyway.recording.output." + outputServiceName
		if !viper.IsSet(outputServiceConfigKey) {
			cobra.CheckErr(fmt.Errorf("output service %s is not configured", outputServiceName))
		}
		outputServiceConfig := viper.Get("skyway.recording.output." + outputServiceName).(map[string]interface{})

		outputService, err := internal.LoadRecordingOutputServiceConfig(outputServiceConfig)
		cobra.CheckErr(err)

		outputService.Service = strings.ToUpper(optionToService[outputServiceName])

		token, err := GenerateAdminToken(appId, secretKey, 3600, []string{})
		cobra.CheckErr(err)

		response, err := internal.CreateSession(channelId, publicationId, contentType, outputService, token, url)
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
	recordingCmd.AddCommand(recordingStartCmd)

	recordingStartCmd.Flags().String("channel-id", "", "Channel ID")
	recordingStartCmd.Flags().String("publication-id", "*", "Publication ID")
	recordingStartCmd.Flags().String("content-type", "", "Content type")
	recordingStartCmd.Flags().String("output-service", "", "Target service of save recording files")
}
