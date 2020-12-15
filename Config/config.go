package Config

import (
	"errors"
	"github.com/spf13/viper"
)

type RSSProxyConfig struct {
	ListenPort string
	HttpProxy string
	RSSInfos map[string]*RSSInfo
}

type RSSInfo struct {
	RSSUrl string
}

func initConfigure() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config") // 设置文件名称（无后缀）
	v.SetConfigType("yaml")   // 设置后缀名 {"1.6以后的版本可以不设置该后缀"}
	v.AddConfigPath(".")      // 设置文件所在路径
	v.Set("verbose", true)    // 设置默认参数

	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.New("error reading config:" + err.Error())
	}

	return v, nil
}

// 读取配置文件，序列化到对应的类
func InitConfigure() (*RSSProxyConfig, error) {
	config, err := initConfigure()
	if err != nil {
		return nil, err
	}

	// 反序列化读取配置到 struct
	var rssProxyConfig RSSProxyConfig
	if err := config.Unmarshal(&rssProxyConfig); err != nil {
		return nil, errors.New("Unmarshal Config File:" + err.Error())
	}

	return &rssProxyConfig, nil
}