package bot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

// VerifyRequest 校验请求是否合法
func VerifyRequest(c *Context) {
	viper.SetConfigName("appConfig")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetDefault("app_secret", "")
	viper.SetDefault("app_key", "")
	err := viper.ReadInConfig()
	if err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatalf("write config failed: %v", err)
		}
	}
	viper.WatchConfig()
	appSecret := viper.GetString("app_secret")
	timestamp := c.Request.Header.Get("timestamp")
	timestampInt64, _ := strconv.ParseInt(timestamp, 10, 64)
	if (time.Now().UnixNano()/1e6 - timestampInt64) > 3600000 {
		log.Println("非法请求！")
		c.Abort()
		return
	}
	sign := c.Request.Header.Get("sign")
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, appSecret)
	hmacCode := hmac.New(sha256.New, []byte(appSecret))
	hmacCode.Write([]byte(stringToSign))
	selfSign := base64.StdEncoding.EncodeToString(hmacCode.Sum(nil))
	if sign != selfSign {
		log.Println("非法请求！")
		c.Abort()
	}
}
