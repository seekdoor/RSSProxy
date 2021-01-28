package NormalRSSHelper

import "github.com/go-resty/resty/v2"

type RSSHelper struct {
	client *resty.Client
}

// 对于一般需要代理的 RSS 使用这个
func NormalRSSHelper(httpProxy string) *RSSHelper {
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

