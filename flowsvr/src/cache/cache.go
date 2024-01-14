package cache

import (
	"asynchub/flowsvr/src/config"
	"time"

	"github.com/niuniumart/gosdk/goredis"
	"github.com/niuniumart/gosdk/martlog"
)

var (
	rdb        *goredis.RedisCli    // redis 连接
	prefix     = "asynchub_kryiea_" // 缓存前缀：
	expireTime = time.Hour * 24     // 过期时间，默认24小时
)

func InitCacche() error {
	goredis.Factory.MaxIdleConn = config.Conf.Redis.MaxIdle                                  // 最大空闲连接数
	goredis.Factory.IdleTimeout = time.Second * time.Duration(config.Conf.Redis.IdleTimeout) // 空闲连接超时时间
	goredis.Factory.MaxConn = config.Conf.Redis.MaxActive

	redisCli, err := goredis.Factory.CreateRedisCli(config.Conf.Redis.Auth, config.Conf.Redis.Url)
	if err != nil {
		martlog.Errorf("redis 连接失败：%v", err)
		return err
	}
	rdb = redisCli

	// 初始化过期时间，默认24小时
	if config.Conf.Redis.CacheTimeoutDay != 0 {
		expireTime = time.Hour * 24 * time.Duration(config.Conf.Redis.CacheTimeoutDay)
	}

	return nil
}
