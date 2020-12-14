package mredis

import (
	"crypto/tls"
	"fmt"
	"github.com/silenceper/wechat/v2/cache"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const WX_ACCESS_TOKEN_KEY = "wx_access_token"

var initOnce sync.Once

type WXRedisCache struct {
	cli *redis.Client
}

var wx_cache *WXRedisCache

func GetRedisCache() cache.Cache {
	initOnce.Do(func() {
		cfg := viper.Get("redis").(map[string]string)
		conn_Addr := fmt.Sprintf("%s:%s", cfg["Host"], cfg["Port"])
		logrus.Infof("%+v", cfg)
		logrus.Info(conn_Addr)
		rcli := redis.NewClient(&redis.Options{
			Addr: conn_Addr,
			DB:   1,
			TLSConfig: &tls.Config{
				ServerName:         cfg["Host"],
				InsecureSkipVerify: true},
		})
		wx_cache = &WXRedisCache{
			cli: rcli,
		}
	})
	return wx_cache
}

func (c *WXRedisCache) Get(key string) interface{} {
	if res, err := c.cli.Get(key).Result(); err != nil {
		logrus.Errorf("Get AccessToken Failed: %v", err.Error())
		return nil
	} else {
		logrus.Infof("Get AccessToken: %s", res)
		return res
	}

}

func (c *WXRedisCache) Set(key string, val interface{}, timeout time.Duration) error {
	logrus.Infof("Set AccessToken: %s:%v", key, val)
	if err := c.cli.Set(key, val, time.Second*7200).Err(); err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (c *WXRedisCache) IsExist(Key string) bool {
	if res, err := c.cli.Exists(Key).Result(); err != nil {
		return false
	} else {
		return res == 1
	}

}

func (c *WXRedisCache) Delete(key string) error {
	return c.cli.Del(key).Err()
}
