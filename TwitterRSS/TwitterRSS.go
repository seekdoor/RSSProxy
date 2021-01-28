package TwitterRSS

import (
	"context"
	"fmt"
	"github.com/gorilla/feeds"
	twitterScraper "github.com/n0madic/twitter-scraper"
	"strings"
	"sync"
)

var (
	// Global mutex
	mu sync.Mutex
)

// Twitter2RSS return RSS from twitter timeline
func Twitter2RSS(httpProxy, screenName string, count int, excludeReplies bool, photoOnly bool) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	feed := &feeds.Feed{
		Title:       "Twitter feed @" + screenName,
		Link:        &feeds.Link{Href: "https://twitter.com/" + screenName},
		Description: "Twitter feed @" + screenName + " through Twitter to RSS proxy by Nomadic",
	}

	err := twitterScraper.SetProxy(httpProxy)
	if err != nil {
		return "", err
	}

	for tweet := range twitterScraper.GetTweets(context.Background(), screenName, count) {
		if tweet.Error != nil {
			return "", tweet.Error
		}

		if excludeReplies && tweet.IsReply {
			continue
		}

		if len(tweet.Photos) <= 0 && photoOnly {
			continue
		}

		if tweet.TimeParsed.After(feed.Created) {
			feed.Created = tweet.TimeParsed
		}

		var title string

		titleSplit := strings.FieldsFunc(tweet.Text, func(r rune) bool {
			return r == '\n' || r == '!' || r == '?' || r == ':' || r == '<' || r == '.' || r == ','
		})
		if len(titleSplit) > 0 {
			if strings.HasPrefix(titleSplit[0], "a href") || strings.HasPrefix(titleSplit[0], "http") {
				title = "link"
			} else {
				title = titleSplit[0]
			}
		}
		title = strings.TrimSuffix(title, "https")
		title = strings.TrimSpace(title)

		content := fmt.Sprintf("<p>%s</p>", tweet.Text)
		for _, oneImageUrl := range tweet.Photos {
			content += fmt.Sprintf("<p><img src=%q></p>", oneImageUrl)
		}
		for _, oneVideoUrl := range tweet.Videos {
			content += fmt.Sprintf("<p><video src=%q></p>", oneVideoUrl)
		}

		feed.Add(&feeds.Item{
			Author:      &feeds.Author{Name: screenName},
			Created:     tweet.TimeParsed,
			Description: content, //tweet.HTML,
			Id:          tweet.PermanentURL,
			Link:        &feeds.Link{Href: tweet.PermanentURL},
			Title:       title,
		})
	}

	if len(feed.Items) == 0 {
		return "", fmt.Errorf("tweets not found")
	}

	return feed.ToRss()
}