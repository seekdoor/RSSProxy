package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/muesli/cache2go"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	var (
		err       error
		vipConfig *viper.Viper
	)
	baseCache = cache2go.Cache("baseCache")
	normalRSSCache = cache2go.Cache("normalRSSCache")
	instagramCache = cache2go.Cache("instagramCache")
	twitterCache = cache2go.Cache("twitterCache")
	// -------------------------------------------------------------
	// 日志设置
	log = NewLogHelper("RSSProxy", logrus.DebugLevel, time.Duration(7*24)*time.Hour, time.Duration(24)*time.Hour)
	// -------------------------------------------------------------
	// 加载配置
	vipConfig, err = InitConfigure()
	if err != nil {
		log.Errorln("InitConfigure:", err)
		return
	}
	// 配置缓存
	ViperConfig2Cache(vipConfig, baseCache, normalRSSCache, instagramCache, twitterCache)
	// 监控配置是否更新
	vipConfig.WatchConfig()
	vipConfig.OnConfigChange(func(e fsnotify.Event) {
		//viper配置发生变化了 执行响应的操作
		log.Infoln("Config file changed:", e.Name)
		// 重新缓存配置
		ViperConfig2Cache(vipConfig, baseCache, normalRSSCache, instagramCache, twitterCache)
		// 立即触发第一次的更新
		getAllSpecialRSS()
		log.Infoln("Update Config Cache Done")
	})
	// -------------------------------------------------------------
	// 初始化 gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 默认所有都通过
	r.Use(cors.Default())
	// 注册 RSSInfos 里面需要中转代理的 RRS 源
	r.GET("rss", routerRSSContent)
	// 注册 IGInfo 里面需要关注的 Instagram 用户
	r.GET("rss/instagram", routerRSSContent4Instagram)
	// 注册 TwiInfo 里面需要关注的 Twitter 用户
	r.GET("rss/twitter", routerRSSContent4Twitter)
	// 手动刷新
	r.GET("refresh", func(context *gin.Context) {
		//viper配置发生变化了 执行响应的操作
		log.Infoln("Manual Refresh Start...")
		// 重新缓存配置
		ViperConfig2Cache(vipConfig, baseCache, normalRSSCache, instagramCache, twitterCache)
		// 立即触发第一次的更新
		getAllSpecialRSS()
		log.Infoln("Manual Refresh Done")
		context.String(http.StatusOK, "Manual Refresh Done")
	})
	// -------------------------------------------------------------
	// 开启一个协程做定时的更新
	// 这里直接使用 cron ，他会自动开启一个协程
	c := cron.New()
	everyTime, err := baseCache.Value("EveryTime")
	if err != nil {
		log.Errorln("baseCache.EveryTime:", err)
		return
	}
	// 定时器
	entryID, err := c.AddFunc("@every " + everyTime.Data().(string), func() {
		getAllSpecialRSS()
	})
	if err != nil {
		log.Errorln("cron entryID:", entryID, "Error:", err)
		return
	}
	// 立即触发第一次的更新
	getAllSpecialRSS()
	c.Start()
	defer c.Stop()
	// -------------------------------------------------------------
	// 启动 gin
	log.Infoln("RSSProxy Start...")
	err = r.Run(":" + vipConfig.GetString("ListenPort"))
	if err != nil {
		log.Fatal("Start RSSProxy Server At Port", vipConfig.GetString("ListenPort"), "Fatal", err)
	}
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


var (
	log       *logrus.Logger
	// 基础的缓存专用
	baseCache  *cache2go.CacheTable
	// 一般的 RSS 的连接配置，比如从 RSSHub 读取的
	normalRSSCache  *cache2go.CacheTable
	// Instagram 的缓存专用，需要这边定时的拉去，然后再提供查询
	instagramCache *cache2go.CacheTable
	// Twitter 的缓存专用，需要这边定时的拉去，然后再提供查询
	twitterCache *cache2go.CacheTable
)