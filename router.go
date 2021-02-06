package main

import (
	"github.com/allanpk716/rssproxy/NormalRSSHelper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func routerRSSContent(c *gin.Context) {
	// like /rss?key=pixiv_month
	nowKey := c.Query("key")
	cacheUrl, err := normalRSSCache.Value(nowKey)
	if err != nil {
		log.Errorln("normalRSSCache." + nowKey +":", err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	nowHttpProxy, err := baseCache.Value("HttpProxy")
	if err != nil {
		log.Errorln("baseCache.HttpProxy:", err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	rssContent, err := NormalRSSHelper.NormalRSSHelper(nowHttpProxy.Data().(string)).GetRSSContent(cacheUrl.Data().(string))
	if err != nil {
		log.Errorln("Get RSS By Proxy Error, RSS Url:", cacheUrl.Data().(string))
		log.Errorln(err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(rssContent))
}

func routerRSSContent4Instagram(c *gin.Context) {
	// like rss/instagram?key=fjamie013
	userYourWantSee := c.Query("key")
	rssContent, err := instagramCache.Value(userYourWantSee)
	if err != nil {
		log.Errorln("instagramCache." + userYourWantSee + ":", err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(rssContent.Data().(string)))
}

func routerRSSContent4Twitter(c *gin.Context) {
	// like rss/twitter?key=baby_eiss
	userYourWantSee := c.Query("key")
	rssContent, err := twitterCache.Value(userYourWantSee)
	if err != nil {
		log.Errorln("twitterCache." + userYourWantSee + ":", err)
		c.JSON(http.StatusNotFound, nil)
		return
	}

	c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(rssContent.Data().(string)))
}