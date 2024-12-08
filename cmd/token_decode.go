package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

func DecodeAdminToken(token string) (SkyWayAdminAuthToken, error) {
	var decoded SkyWayAdminAuthToken
	_, err := jwt.ParseWithClaims(token, &decoded, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})

	return decoded, err
}

func DecodeToken(token string) (SkyWayAuthToken, error) {
	var decoded SkyWayAuthToken
	_, err := jwt.ParseWithClaims(token, &decoded, func(t *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	return decoded, err
}

// tokenDecodeCmd represents the decode command
var tokenDecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode SkyWay Auth Token or SkyWay Admin Auth Token",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(cmd.InOrStdin())
		stdinBytes, err := reader.ReadBytes('\n')
		cobra.CheckErr(err)

		isAdmin, err := cmd.Flags().GetBool("admin")
		cobra.CheckErr(err)

		pretty, err := cmd.Flags().GetBool("pretty")
		cobra.CheckErr(err)

		if isAdmin {
			decoded, err := DecodeAdminToken(string(stdinBytes))
			// Errors are ignored because the token is only decoded, not verified.
			if err != nil && !errors.Is(err, jwt.ErrSignatureInvalid) {
				cobra.CheckErr(err)
			}

			if pretty {
				token, err := json.MarshalIndent(decoded, "", "  ")
				cobra.CheckErr(err)

				fmt.Println(string(token))
			} else {
				token, err := json.Marshal(decoded)
				cobra.CheckErr(err)

				fmt.Println(string(token))
			}
			return
		} else {
			decoded, err := DecodeToken(string(stdinBytes))
			// Errors are ignored because the token is only decoded, not verified.
			if err != nil && !errors.Is(err, jwt.ErrSignatureInvalid) {
				cobra.CheckErr(err)
			}

			if pretty {
				token, err := json.MarshalIndent(decoded, "", "  ")
				cobra.CheckErr(err)

				fmt.Println(string(token))
			} else {
				token, err := json.Marshal(decoded)
				cobra.CheckErr(err)

				fmt.Println(string(token))
			}
		}
	},
}

func init() {
	tokenCmd.AddCommand(tokenDecodeCmd)

	tokenDecodeCmd.Flags().BoolP("admin", "a", false, "Decode SkyWay Admin Auth Token")
	tokenDecodeCmd.Flags().BoolP("pretty", "p", false, "Pretty print JSON")
}
