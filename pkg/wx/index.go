package wx

import (
	"sync"
	"we_server/pkg/redis"

	"github.com/silenceper/wechat/v2"
	"github.com/spf13/viper"

	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

var account *officialaccount.OfficialAccount
var initOnce sync.Once

func GetAccount() *officialaccount.OfficialAccount {
	initOnce.Do(func() {
		wc := wechat.NewWechat()
		cfg := offConfig.Config{
			AppID:     viper.GetString("WX_APPID"),
			AppSecret: viper.GetString("WX_APPSECRET"),
			Cache:     mredis.GetRedisCache(),
			Token:     viper.GetString("WX_APPTOKEN"),
		}
		account = wc.GetOfficialAccount(&cfg)
	})
	return account
}
