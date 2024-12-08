package cmd_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kadoshita/skyway-cli/cmd"
)

var appId = "402e60fb-9698-4eb9-9ee2-6d3d66a78068"
var secretKey = "5ecccc3b-5577-4747-b04e-df3194f54b63"

func TestVerify(t *testing.T) {
	t.Run("verify", func(t *testing.T) {
		t.Run("正常なトークンの場合はエラーが発生しない", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err != nil {
				t.Error("正常なトークンの場合はエラーが発生しない")
			}
		})
		t.Run("シークレットキーが異なる場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(uuid.New().String()))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("シークレットキーが異なる場合はエラー")
			}
		})
		t.Run("ペイロードがJSON形式ではない場合はエラー", func(t *testing.T) {
			header := `{"alg":"HS256","typ":"JWT"}`
			encodedHeader := base64.RawURLEncoding.EncodeToString([]byte(header))
			encodedPayload := base64.RawURLEncoding.EncodeToString([]byte("{invalid-payload}"))
			tokenWithoutSignature := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)
			hmac := hmac.New(sha256.New, []byte(secretKey))
			hmac.Write([]byte(tokenWithoutSignature))
			signature := strings.TrimRight(base64.RawURLEncoding.EncodeToString(hmac.Sum(nil)), "=")
			tokenString := fmt.Sprintf("%s.%s", tokenWithoutSignature, signature)

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("ペイロードがJSON形式ではない場合はエラー")
			}
		})
	})
	t.Run("jti", func(t *testing.T) {
		t.Run("jtiが無い場合はエラー", func(t *testing.T) {
			claims := jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}

			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("jtiが無い場合はエラー")
			}
		})
		t.Run("jtiがUUID v4形式ではない場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     "invalid-jti",
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("jtiがUUID v4形式ではない場合はエラー")
			}
		})
	})
	t.Run("iat", func(t *testing.T) {
		t.Run("iatが無い場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("iatが無い場合はエラー")
			}
		})
		t.Run("iatが秒単位のUNIXタイムスタンプではない場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().UnixNano()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("iatが秒単位のUNIXタイムスタンプではない場合はエラー")
			}
		})
		t.Run("iatが現在時刻の2分後よりも後の場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Add(time.Duration(2)*time.Minute + time.Duration(2)*time.Second).Unix()), // マージンとして2秒追加
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("iatが現在時刻の2分後よりも後の場合はエラー")
			}
		})
	})
	t.Run("exp", func(t *testing.T) {
		t.Run("expが無い場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("expが無い場合はエラー")
			}
		})
		t.Run("expが秒単位のUNIXタイムスタンプではない場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).UnixNano()), // ナノ秒単位
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("expが秒単位のUNIXタイムスタンプではない場合はエラー")
			}
		})
		t.Run("expがiatから72時間を超えている場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(72)*time.Hour + time.Duration(1)*time.Second).Unix()), // 72時間1秒
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("expがiatから72時間を超えている場合はエラー")
			}
		})
	})
	t.Run("version", func(t *testing.T) {
		t.Run("versionが未指定、1、2以外の場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 3,
				"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("versionが未指定、1、2以外の場合はエラー")
			}
		})
	})
	t.Run("scope", func(t *testing.T) {
		t.Run("scopeが無い場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scopeが無い場合はエラー")
			}
		})
		t.Run("scopeが数値の場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(72)*time.Hour + time.Duration(1)*time.Second).Unix()), // 72時間1秒
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   42,
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scopeが数値の場合はエラー")
			}
		})
		t.Run("scopeが文字列の場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(72)*time.Hour + time.Duration(1)*time.Second).Unix()), // 72時間1秒
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   "invalid-scope",
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scopeが文字列の場合はエラー")
			}
		})
		t.Run("scopeが配列の場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(72)*time.Hour + time.Duration(1)*time.Second).Unix()), // 72時間1秒
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   []int{},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scopeが配列の場合はエラー")
			}
		})
	})
	t.Run("scope.app", func(t *testing.T) {
		t.Run("scope.appが無い場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scope.appが無い場合はエラー")
			}
		})
		t.Run("scope.appがオブジェクトでは無い場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": []int{}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scope.appがオブジェクトでは無い場合はエラー")
			}
		})
	})
	t.Run("scope.app.id", func(t *testing.T) {
		t.Run("scope.app.idが無い場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope":   map[string]interface{}{"app": map[string]interface{}{}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scope.app.idが無い場合はエラー")
			}
		})
		t.Run("scope.app.idがUUID v4形式ではない場合はエラー", func(t *testing.T) {
			tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat":     int(time.Now().Unix()),
				"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
				"jti":     uuid.New().String(),
				"version": 2,
				"scope": map[string]interface{}{"app": map[string]interface{}{
					"id": "invalid-app-id",
				}},
			}).SignedString([]byte(secretKey))

			err := cmd.VerifyToken(tokenString, secretKey)

			if err == nil {
				t.Error("scope.app.idがUUID v4形式ではない場合はエラー")
			}
		})
	})
	t.Run("tokenのサイズが7KBを超えている場合はエラー", func(t *testing.T) {
		// 7KBのランダムな文字列を生成
		sevenKB := make([]byte, 7*1024)
		for i := range sevenKB {
			sevenKB[i] = 'a'
		}
		claims := jwt.MapClaims{
			"iat":     int(time.Now().Unix()),
			"exp":     int(time.Now().Add(time.Duration(60) * time.Second).Unix()),
			"jti":     uuid.New().String(),
			"version": 2,
			"scope":   map[string]interface{}{"app": map[string]interface{}{"id": appId, "actions": []string{"read", string(sevenKB)}}},
		}
		tokenString, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

		err := cmd.VerifyToken(tokenString, secretKey)

		if err == nil {
			t.Error("tokenのサイズが7KBを超えている場合はエラー")
		}
	})
}
