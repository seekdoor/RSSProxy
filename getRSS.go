package main

import (
	"homelab/rssproxy/InstagramRSS"
	"homelab/rssproxy/TwitterRSS"
	"os"
	"sync"
)

var (
	// 后台更新 Instagram 的时候，只允许一个运行，所以需要加锁
	mutexInstagram         sync.Mutex
	updateInstagramRunning bool

	// 后台更新 Twitter     的时候，只允许一个运行，所以需要加锁
	mutexTwitter        	sync.Mutex
	updateTwitterRunning 	bool
)

// 统一一起获取
func getAllSpecialRSS()  {
	go getInstagramRSS()
	go getTwitterRSS()
}

func getInstagramRSS() {
	mutexInstagram.Lock()
	log.Infoln("------- Start Update Instagram RSS -------")
	defer func() {
		updateInstagramRunning = false
		mutexInstagram.Unlock()
		log.Infoln("------- End Update Instagram RSS -------")
	}()

	if updateInstagramRunning == true {
		log.Infoln("getInstagramRSS is running, so skip this time")
		return
	}
	// 设置已经在后台刷新
	updateInstagramRunning = true

	userName, err := instagramCache.Value("IGInfo.UserName")
	if err != nil {
		log.Errorln("instagramCache:IGInfo.UserName:", err)
		return
	}
	passWord, err := instagramCache.Value("IGInfo.PassWord")
	if err != nil {
		log.Errorln("instagramCache.IGInfo.PassWord:", err)
		return
	}
	httpProxy, err := baseCache.Value("HttpProxy")
	if err != nil {
		log.Errorln("baseCache.HttpProxy:", err)
		return
	}
	feedMaxItems, err := instagramCache.Value("IGInfo.FeedMaxItems")
	if err != nil {
		log.Errorln("instagramCache.IGInfo.FeedMaxItems:", err)
		return
	}
	instagramUsers, err := instagramCache.Value("IGInfo.InstagramUsers")
	if err != nil {
		log.Errorln("instagramCache.IGInfo.InstagramUsers:", err)
		return
	}

	rssMap, err := InstagramRSS.NewInstagramRSS(userName.Data().(string), passWord.Data().(string),
		httpProxy.Data().(string), feedMaxItems.Data().(int), instagramUsers.Data().([]string))
	if err != nil {
		log.Errorln("Get InstagramRSS.InstagramRSS Error:", err)

		// 删除缓存的登录文件
		if Exists("igCache") {
			if err := os.Remove("igCache"); err != nil {
				log.Errorln("Del igCache:", err)
			}
		}

		return
	}

	for k, v := range rssMap {
		instagramCache.Add(k, 0 , v)
	}
}

func getTwitterRSS() {
	mutexTwitter.Lock()
	log.Infoln("------- Start Update Twitter RSS -------")
	defer func() {
		updateTwitterRunning = false
		mutexTwitter.Unlock()
		log.Infoln("------- End Update Twitter RSS -------")
	}()

	if updateTwitterRunning == true {
		log.Infoln("getTwitterRSS is running, so skip this time")
		return
	}
	// 设置已经在后台刷新
	updateTwitterRunning = true

	httpProxy, err := baseCache.Value("HttpProxy")
	if err != nil {
		log.Errorln("baseCache.HttpProxy:", err)
		return
	}
	feedMaxItems, err := twitterCache.Value("TwiInfo.FeedMaxItems")
	if err != nil {
		log.Errorln("twitterCache.TwiInfo.FeedMaxItems:", err)
		return
	}

	excludeReplies, err := twitterCache.Value("TwiInfo.ExcludeReplies")
	if err != nil {
		log.Errorln("twitterCache:TwiInfo.ExcludeReplies:", err)
		return
	}

	photoOnly, err := twitterCache.Value("TwiInfo.PhotoOnly")
	if err != nil {
		log.Errorln("twitterCache:TwiInfo.PhotoOnly:", err)
		return
	}

	twitterUsers, err := twitterCache.Value("TwiInfo.TwitterUsers")
	if err != nil {
		log.Errorln("twitterCache.TwiInfo.TwitterUsers:", err)
		return
	}

	for _, oneUserID := range twitterUsers.Data().([]string) {
		tmpRSSString, err := TwitterRSS.Twitter2RSS(httpProxy.Data().(string), oneUserID, feedMaxItems.Data().(int), excludeReplies.Data().(bool), photoOnly.Data().(bool))
		if err != nil {
			log.Errorln("Get TwitterRSS.Twitter2RSS", oneUserID, "error", err.Error())
			continue
		}

		twitterCache.Add(oneUserID, 0, tmpRSSString)
		log.Infoln("Update Twitter:", oneUserID, "Done")
	}

}

