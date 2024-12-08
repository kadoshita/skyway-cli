package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

type SkyWayAdminAuthToken struct {
	Jti           string `json:"jti"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
	AppId         string `json:"appId"`
	jwt.MapClaims `json:"-"`
}

type SkyWayAuthToken struct {
	Jti           string               `json:"jti"`
	Iat           int                  `json:"iat"`
	Exp           int                  `json:"exp"`
	Version       int                  `json:"version"`
	Scope         SkyWayAuthTokenScope `json:"scope"`
	jwt.MapClaims `json:"-"`
}
type SkyWayAuthTokenScope struct {
	App SkyWayAuthTokenAppScope `json:"app"`
}
type SkyWayAuthTokenAppScope struct {
	Id        string                        `json:"id"`
	Actions   []string                      `json:"actions"`
	Turn      bool                          `json:"turn"`
	Analytics bool                          `json:"analytics"`
	Channels  []SkyWayAuthTokenChannelScope `json:"channels"`
}
type SkyWayAuthTokenChannelScope struct {
	Id      string                       `json:"id"`
	Name    string                       `json:"name"`
	Actions []string                     `json:"actions"`
	Members []SkyWayAuthTokenMemberScope `json:"members"`
	SfuBots []SkyWayAuthTokenSfuBotScope `json:"sfuBots"`
}
type SkyWayAuthTokenMemberScope struct {
	Id           string                           `json:"id"`
	Name         string                           `json:"name"`
	Actions      []string                         `json:"actions"`
	Publication  SkyWayAuthTokenPublicationScope  `json:"publication"`
	Subscription SkyWayAuthTokenSubscriptionScope `json:"subscription"`
}
type SkyWayAuthTokenPublicationScope struct {
	Actions []string `json:"actions"`
}
type SkyWayAuthTokenSubscriptionScope struct {
	Actions []string `json:"actions"`
}
type SkyWayAuthTokenSfuBotScope struct {
	Actions     []string                         `json:"actions"`
	Forwardings []SkyWayAuthTokenForwardingScope `json:"forwardings"`
}
type SkyWayAuthTokenForwardingScope struct {
	Actions []string `json:"actions"`
}

func modifyTokenTemplate(tokenTmpl string, values []string) (string, error) {
	for _, v := range values {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			return "", fmt.Errorf("invalid value: %s", v)
		}
		modifiedTokenTmpl, err := sjson.Set(tokenTmpl, kv[0], kv[1])
		if err != nil {
			return "", fmt.Errorf("failed to set value. key: %s, value: %s", kv[0], kv[1])
		}
		tokenTmpl = modifiedTokenTmpl
		slog.Debug("Set value", "key", kv[0], "value", kv[1])
	}
	return tokenTmpl, nil
}

func GenerateToken(tokenTmpl string, appId string, secretKey string, expire int, values []string) (string, error) {
	var token SkyWayAuthToken

	modifiedTokenTmpl, err := modifyTokenTemplate(tokenTmpl, []string{})
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal([]byte(modifiedTokenTmpl), &token); err != nil {
		return "", fmt.Errorf("invalid token template. skyway.token.tmpl: %s", tokenTmpl)
	}

	// jti, iat, exp and appId are overwrited by command line arguments
	token.Jti = uuid.New().String()
	token.Iat = int(time.Now().Unix())
	token.Exp = int(time.Now().Add(time.Duration(expire) * time.Second).Unix())
	token.Scope.App.Id = appId

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token. err: %v", err)
	}
	return tokenString, nil
}

func GenerateAdminToken(appId string, secretKey string, expire int, values []string) (string, error) {
	var token SkyWayAdminAuthToken
	token.Jti = uuid.New().String()
	token.Iat = int(time.Now().Unix())
	token.Exp = int(time.Now().Add(time.Duration(expire) * time.Second).Unix())
	token.AppId = appId

	for _, v := range values {
		kv := strings.Split(v, "=")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "jti":
			token.Jti = kv[1]
		case "iat":
			var err error
			token.Iat, err = strconv.Atoi(kv[1])
			if err != nil {
				return "", fmt.Errorf("failed to convert iat. iat: %s", kv[1])
			}
		case "exp":
			var err error
			token.Exp, err = strconv.Atoi(kv[1])
			if err != nil {
				return "", fmt.Errorf("failed to convert exp. exp: %s", kv[1])
			}
		case "appId":
			token.AppId = kv[1]
		}
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, token).SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token. err: %v", err)
	}
	return tokenString, nil
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "SkyWay Auth Token Generate Decode and Verify",
	Long: `Generate, decode and verify SkyWay Auth Token.
when do not specify subcommand, generate token.
Note: This command requires the skyway.token.tmpl configuration for generate token.
This command cannot append additional scopes to token template configuration yet.`,
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

		isAdmin, err := cmd.Flags().GetBool("admin")
		cobra.CheckErr(err)

		if isAdmin {
			tokenString, err := GenerateAdminToken(appId, secretKey, expire, []string{})
			cobra.CheckErr(err)

			fmt.Println(tokenString)
		} else {
			tokenTmpl := viper.GetString("skyway.token.tmpl")

			tokenString, err := GenerateToken(tokenTmpl, appId, secretKey, expire, []string{})
			cobra.CheckErr(err)

			fmt.Println(tokenString)
		}
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)

	tokenCmd.Flags().String("app-id", "", "SkyWay App ID. This option can also be set by the skyway.app_id configuration or the SKYWAY_APP_ID environment variable.")
	tokenCmd.Flags().String("secret-key", "", "SkyWay Secret Key. This option can also be set by the skyway.secret_key configuration or the SKYWAY_SECRET_KEY environment variable.")
	tokenCmd.Flags().Int("expire", 3600, "Token expire time in seconds. This option can also be set by the skyway.token.expire configuration or the SKYWAY_TOKEN_EXPIRE environment variable.")
	tokenCmd.Flags().String("tmpl", "", "Token template. This option can also be set by the skyway.token.tmpl configuration or the SKYWAY_TOKEN_TMPL environment variable.")

	tokenCmd.Flags().Bool("admin", false, "Generate SkyWay Admin Auth Token")
}
