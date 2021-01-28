package main

import (
	"github.com/muesli/cache2go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func InitConfigure() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config") // 设置文件名称（无后缀）
	v.SetConfigType("yaml")   // 设置后缀名 {"1.6以后的版本可以不设置该后缀"}
	v.AddConfigPath(".")      // 设置文件所在路径

	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.New("error reading config:" + err.Error())
	}

	return v, nil
}

func Add2CacheString(config *viper.Viper, cache *cache2go.CacheTable, keyName string) {
	cache.Add(keyName, 0, config.GetString(keyName))
}

func Add2CacheInt(config *viper.Viper, cache *cache2go.CacheTable, keyName string) {
	cache.Add(keyName, 0, config.GetInt(keyName))
}

func Add2CacheBool(config *viper.Viper, cache *cache2go.CacheTable, keyName string) {
	cache.Add(keyName, 0, config.GetBool(keyName))
}

func ViperConfig2Cache(config *viper.Viper, cacheBase, rsscahe, igcache, twicache *cache2go.CacheTable) {
	// ------------------------------------------------------------
	// 基础配置信息
	Add2CacheString(config, cacheBase, "ListenPort")
	Add2CacheString(config, cacheBase, "HttpProxy")
	Add2CacheString(config, cacheBase, "EveryTime")

	// ------------------------------------------------------------
	// 读取 RSSHub 相关需要订阅的信息
	rsscahe.Flush()
	rsshubInfo := config.GetStringMapString("RSSInfos")
	for k, v := range rsshubInfo {
		rsscahe.Add(k, 0, v)
		println("Support Router Key:", "rss/?key=" + k, "--", v)
	}
	// ------------------------------------------------------------
	// 读取 IGInfo 相关配置
	Add2CacheString(config, igcache, "IGInfo.UserName")
	Add2CacheString(config, igcache, "IGInfo.PassWord")
	Add2CacheInt(config, igcache, "IGInfo.FeedMaxItems")
	IGUsers := config.GetStringSlice("IGInfo.InstagramUsers")
	for _, oneUser := range IGUsers{
		println("Support Router Key:", "rss/instagram/?key=" + oneUser, "--", "https://www.instagram.com/" + oneUser)
	}
	igcache.Add("IGInfo.InstagramUsers", 0, IGUsers)
	// ------------------------------------------------------------
	// 读取 TwiInfo 相关配置
	Add2CacheInt(config, twicache, "TwiInfo.FeedMaxItems")
	Add2CacheBool(config, twicache,"TwiInfo.ExcludeReplies")
	Add2CacheBool(config, twicache,"TwiInfo.PhotoOnly")
	TwiUsers := config.GetStringSlice("TwiInfo.TwitterUsers")
	for _, oneUser := range TwiUsers{
		println("Support Router Key:", "rss/twitter/?key=" + oneUser, "--", "https://www.twitter.com/" + oneUser)
	}
	twicache.Add("TwiInfo.TwitterUsers", 0, TwiUsers)
	// ------------------------------------------------------------
}
