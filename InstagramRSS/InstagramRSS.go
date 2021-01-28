package InstagramRSS

import (
	"fmt"
	"github.com/Masterminds/goutils"
	"github.com/ahmdrz/goinsta/v2"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
	"time"
)

var(
	insta        *goinsta.Instagram
	feedMaxItems int
)

func NewInstagramRSS(igUserName, igPassword, httpProxy string, inFeedMaxItems int, instagramUsers []string) (map[string]string, error) {
	var (
		configFile string
		igUsers []string
		userRSS map[string]string
		err     error
	)

	userRSS = map[string]string{}

	igUsers = instagramUsers
	feedMaxItems = inFeedMaxItems
	configFile = "igCache"

	if insta, err = goinsta.Import(configFile); err != nil {
		println("Unable to import Instagram configuration:", err)
		println("Attempting new login")

		insta = goinsta.New(igUserName, igPassword)
		if err = insta.SetProxy(httpProxy, true); err != nil {
			return userRSS, errors.Errorf("SetProxy Error:", err)
		}
		if err = insta.Login(); err != nil {
			return userRSS, errors.Errorf("unable to initialize Instagram client:", err)
		}
	}

	if err := insta.Export(configFile); err != nil {
		return userRSS, errors.Errorf("error: unable to export Instagram client configuration:", err)
	}
	if err = insta.SetProxy(httpProxy, true); err != nil {
		return userRSS, errors.Errorf("SetProxy Error:", err)
	}

	if len(igUsers) == 0 {
		// If no static list of IG users is provided, attempt retrieving the list of followings
		// from the logged user's account
		for followings := insta.Account.Following(); followings.Next(); {
			for _, u := range followings.Users {
				igUsers = append(igUsers, u.Username)
			}

			if err := followings.Error(); err != nil {
				if err == goinsta.ErrNoMore {
					break
				}

				return userRSS, errors.Errorf("unable to retrieve followings:", err)
			}
		}

		if len(igUsers) == 0 {
			return userRSS, errors.Errorf("no users provided")
		}
	}

	for _, u := range igUsers {

		feed := &feeds.Feed{
			Title:       fmt.Sprintf("Instagram %s", u),
			Link:        &feeds.Link{Href: "https://www.instagram.com/" + u},
			Description: "Instagram RSS feed generated RSSProxy",
			Created:     time.Now(),
		}

		items, err := fetchUserFeedItems(u)
		if err != nil {
			return userRSS, errors.Errorf("unable to retrieve user", u, "feed:", err)
		}

		for _, item := range items {
			feed.Add(item)
		}

		feed.Sort(func(a, b *feeds.Item) bool {
			return a.Created.After(b.Created)
		})

		rss, err := feed.ToRss()
		if err != nil {
			return userRSS, errors.Errorf("unable to render RSS feed: %s", err)
		}

		userRSS[u] = rss
	}

	return userRSS, nil
}

func fetchUserFeedItems(name string) ([]*feeds.Item, error) {
	items := make([]*feeds.Item, 0)

	user, err := insta.Profiles.ByName(name)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get user information")
	}

	latest := user.Feed()
	if latest == nil {
		return nil, errors.Wrap(err, "unable to get user latest feed")
	}

	for latest.Next(false) {
		for _, item := range latest.Items {
			items = append(items, formatFeedItem(&item))
			if feedMaxItems != 0 {
				if len(items) >= feedMaxItems {
					return items, nil
				}
			}
		}

		if err := latest.Error(); err != nil {
			if err := latest.Error(); err == goinsta.ErrNoMore {
				break
			}
			return nil, errors.Wrap(err, "unable to retrieve user feed")
		}
	}

	return items, nil
}

func formatFeedItem(item *goinsta.Item) *feeds.Item {
	shortDesc, _ := goutils.Abbreviate(item.Caption.Text, 50)

	content := fmt.Sprintf("<p>%s</p>", item.Caption.Text)

	if len(item.Images.Versions) > 0 {
		content += fmt.Sprintf("<p><img src=%q></p>", item.Images.Versions[0].URL)
	}

	if len(item.CarouselMedia) > 0 {
		for _, i := range item.CarouselMedia {
			content += fmt.Sprintf("<p><img src=%q></p>", i.Images.Versions[0].URL)
		}
	}

	return &feeds.Item{
		Id:      item.ID,
		Title:   shortDesc,
		Created: time.Unix(item.TakenAt, 0),
		Author: &feeds.Author{
			Name:  fmt.Sprintf("%s (@%s)", item.User.FullName, item.User.Username),
			Email: item.User.Username + "@instagram.com",
		},
		Description: shortDesc,
		Content:     content,
		Link:        &feeds.Link{Href: fmt.Sprintf("https://www.instagram.com/p/%s", item.Code)},
	}
}
