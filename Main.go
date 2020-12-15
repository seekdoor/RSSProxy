package main

import (
	"HomeLab/RSSProxy/Config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"io"
	"os"
	"path"
	"sync"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"net/http"
	"strings"
	"time"
)

func main() {
	var err error
	// 日志设置
	log = NewLogHelper("RSSProxy", logrus.DebugLevel, time.Duration(7*24)*time.Hour, time.Duration(24)*time.Hour)
	// 加载配置
	config, err = Config.InitConfigure()
	if err != nil {
		log.Errorln(err)
		return
	}

	// 初始化 gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 默认所有都通过
	r.Use(cors.Default())

	for key, rssinfo := range config.RSSInfos {
		// 分割 key 作为路由
		stringSlice := strings.Split(key,"_")
		if len(stringSlice) < 2 {
			log.Errorln("RssInfo Key Error:", key)
			return
		}
		// 存储到并发专用读取的字典中
		rssKeyMap.Store(key, rssinfo.RSSUrl)

		var routerName string
		for index, oneKey := range stringSlice {
			if index == 0 {
				routerName = "/" + oneKey + "/"
				continue
			}
			if index == len(stringSlice) - 1 {
				routerName += oneKey
			} else {
				routerName += oneKey + "_"
			}

		}
		r.GET(routerName, getRSSContent)
	}

	// 启动 gin
	err = r.Run(":" + config.ListenPort)
	if err != nil {
		log.Fatal("Start RSSProxy Server At Port", config.ListenPort, "Fatal", err)
	}
}

func getRSSContent(c *gin.Context) {
	// like /pixiv/month
	rssKey := strings.Replace(c.FullPath(), "/", "", 1)
	rssKey = strings.Replace(rssKey, "/", "_", -1)

	urlNeedGet, bok := rssKeyMap.Load(rssKey)
	if bok == false {
		log.Errorln("Load rssKeyMap False, key:", rssKey)
		c.JSON(http.StatusNotFound, nil)
		return
	}
	rssContent, err := NewRSSHelper(config.HttpProxy).GetRSSContent(urlNeedGet.(string))
	if err != nil {
		log.Errorln("Get RSS By Proxy Error, RSS Url:", urlNeedGet.(string))
		log.Errorln(err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(rssContent))
}

func NewLogHelper(appName string, level logrus.Level, maxAge time.Duration, rotationTime time.Duration) *logrus.Logger {

	Logger := &logrus.Logger{
		// Out:   os.Stderr,
		// Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}
	nowpath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pathRoot := path.Join(nowpath, "Logs")
	fileAbsPath := path.Join(pathRoot, appName+".log")
	// 下面配置日志每隔 X 分钟轮转一个新文件，保留最近 X 分钟的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		path.Join(pathRoot, appName+"--%Y%m%d%H%M--.log"),
		rotatelogs.WithLinkName(fileAbsPath),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotationTime),
	)

	Logger.SetLevel(level)
	Logger.SetOutput(io.MultiWriter(os.Stderr, writer))

	return Logger
}

type RSSHelper struct {
	client *resty.Client
}

func NewRSSHelper(httpProxy string) *RSSHelper {
	rsshelper := RSSHelper{}
	// Create a Resty Client
	rsshelper.client = resty.New()
	// Setting a Proxy URL and Port
	rsshelper.client.SetProxy(httpProxy)

	return &rsshelper
}

func (r RSSHelper) GetRSSContent(url string) (string, error) {
	resp, err := r.client.R().Get(url)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

var (
	log *logrus.Logger
	config *Config.RSSProxyConfig
	rssKeyMap sync.Map
)