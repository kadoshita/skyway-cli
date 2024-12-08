package cmd

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tokenServeCmd represents the serve command
var tokenServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve SkyWay Auth Token by HTTP Server",
	Long:  `Serve SkyWay Auth Token when send GET request to /token.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.app_id", cmd.Flags().Lookup("app-id"))
		viper.BindPFlag("skyway.secret_key", cmd.Flags().Lookup("secret-key"))
		viper.BindPFlag("skyway.token.expire", cmd.Flags().Lookup("expire"))
		viper.BindPFlag("skyway.token.tmpl", cmd.Flags().Lookup("tmpl"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		appId := viper.GetString("skyway.app_id")
		secretKey := viper.GetString("skyway.secret_key")
		expire := viper.GetInt("skyway.token.expire")
		tokenTmpl := viper.GetString("skyway.token.tmpl")
		port, err := cmd.Flags().GetInt("port")
		cobra.CheckErr(err)

		server := echo.New()
		server.HideBanner = true
		server.Use(middleware.Logger())

		server.GET("/token", func(c echo.Context) error {
			tokenString, err := GenerateToken(tokenTmpl, appId, secretKey, expire, []string{})
			if err != nil {
				return c.String(500, err.Error())
			}
			return c.String(200, tokenString)
		})
		server.Start(":" + fmt.Sprint(port))
	},
}

func init() {
	tokenCmd.AddCommand(tokenServeCmd)

	tokenServeCmd.Flags().String("app-id", "", "SkyWay App ID. This option can also be set by the skyway.app_id configuration or the SKYWAY_APP_ID environment variable.")
	tokenServeCmd.Flags().String("secret-key", "", "SkyWay Secret Key. This option can also be set by the skyway.secret_key configuration or the SKYWAY_SECRET_KEY environment variable.")
	tokenServeCmd.Flags().Int("expire", 3600, "Token expire time in seconds. This option can also be set by the skyway.token.expire configuration or the SKYWAY_TOKEN_EXPIRE environment variable.")
	tokenServeCmd.Flags().String("tmpl", "", "Token template. This option can also be set by the skyway.token.tmpl configuration or the SKYWAY_TOKEN_TMPL environment variable.")

	tokenServeCmd.Flags().IntP("port", "p", 8080, "HTTP Server Port")
}
