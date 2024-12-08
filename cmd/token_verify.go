package cmd

import (
	"bufio"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VerifyToken(tokenString string, secretKey string) error {
	var decoded = jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, decoded, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	var jti string
	var iat int64
	var exp int64
	var version int
	var scope interface{}
	for key, value := range decoded {
		if key == "jti" {
			jti = value.(string)
		}
		if key == "iat" {
			iat = int64(value.(float64))
		}
		if key == "exp" {
			exp = int64(value.(float64))
		}
		if key == "scope" {
			scope = value
		}
		if key == "version" {
			version = int(value.(float64))
		}
	}

	if jti == "" {
		return fmt.Errorf("jti is required")
	}
	parsed, err := uuid.Parse(jti)
	if err != nil {
		return fmt.Errorf("jti should be UUID. value: %s", jti)
	}
	if parsed.Version() != 4 {
		return fmt.Errorf("jti should be UUID v4. value: %s version: %d", jti, parsed.Version())
	}

	if iat == 0 {
		return fmt.Errorf("iat is required")
	}
	// 秒単位のunixタイムスタンプの桁数は10桁なので、iatが10桁ではない場合はエラー
	if len(strconv.FormatInt(iat, 10)) != 10 {
		return fmt.Errorf("iat should be seconds unix timestamp. value: %d", iat)
	}
	// iatが現在時刻の2分後よりも後の場合はエラー
	if iat > (int64(time.Now().Add(2 * time.Minute).Unix())) {
		return fmt.Errorf("iat should be less than 2 minutes from now. value: %d", iat)
	}

	if exp == 0 {
		return fmt.Errorf("exp is required")
	}
	// 秒単位のunixタイムスタンプの桁数は10桁なので、expが10桁ではない場合はエラー
	if len(strconv.FormatInt(exp, 10)) != 10 {
		return fmt.Errorf("exp should be seconds unix timestamp. value: %d", exp)
	}
	// expがiatから72時間を超えている場合はエラー
	if exp > (iat + 60*60*72) {
		return fmt.Errorf("exp should be less than 72 hours from iat. iat: %d exp: %d", iat, exp)
	}

	if version != 0 && version != 1 && version != 2 {
		return fmt.Errorf("version should be undefined, 1 or 2. value: %d", version)
	}

	if scope == nil {
		return fmt.Errorf("scope is required")
	}
	// scopeがオブジェクトでは無い場合はエラー
	if _, ok := scope.(map[string]interface{}); !ok {
		return fmt.Errorf("scope should be object. value: %v", scope)
	}

	scopeMap := scope.(map[string]interface{})
	if _, ok := scopeMap["app"]; !ok {
		return fmt.Errorf("scope.app is required")
	}
	// scope.appがオブジェクトでは無い場合はエラー
	if _, ok := scopeMap["app"].(map[string]interface{}); !ok {
		return fmt.Errorf("scope.app should be object. value: %v", scopeMap["app"])
	}

	appMap := scopeMap["app"].(map[string]interface{})
	if _, ok := appMap["id"]; !ok {
		return fmt.Errorf("scope.app.id is required")
	}
	// scope.app.idがUUID v4形式ではない場合はエラー
	appId := appMap["id"].(string)
	parsed, err = uuid.Parse(appId)
	if err != nil {
		return fmt.Errorf("scope.app.id should be UUID. value: %s", appId)
	}
	if parsed.Version() != 4 {
		return fmt.Errorf("scope.app.id should be UUID v4. value: %s version: %d", appId, parsed.Version())
	}

	// tokenのサイズが7KBを超えている場合はエラー
	if len(tokenString) > 7*1024 {
		return fmt.Errorf("token size should be less than 7KB. size: %d", len(tokenString))
	}

	return nil
}

// tokenVerifyCmd represents the verify command
var tokenVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify SkyWay Auth Token",
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("skyway.secret_key", cmd.Flags().Lookup("secret-key"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(cmd.InOrStdin())
		stdinBytes, err := reader.ReadBytes('\n')
		cobra.CheckErr(err)

		secretKey := viper.GetString("skyway.secret_key")

		err = VerifyToken(string(stdinBytes), secretKey)
		cobra.CheckErr(err)

		fmt.Println("Token is valid")
	},
}

func init() {
	tokenCmd.AddCommand(tokenVerifyCmd)

	tokenVerifyCmd.Flags().String("secret-key", "", "SkyWay Secret Key. This option can also be set by the skyway.secret_key configuration or the SKYWAY_SECRET_KEY environment variable.")
}
