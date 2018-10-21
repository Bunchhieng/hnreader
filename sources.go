package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jzelinskie/geddit"
)

// SteemItSource fetches latest news from SteemIt
type SteemItSource struct{}

// Fetch gets data from SteemIt
func (l *SteemItSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	resp, err := http.Get(SteemItURL)
	if err != nil {
		return news, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return news, err
	}

	m := make(map[string]map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return news, err
	}

	idx := 0
	for _, value := range m {
		if idx >= count {
			break
		}

		str, ok := value["url"].(string)
		if ok {
			news[idx] = str
			idx++
		}
	}

	return news, nil
}

// DevToSource fetches latest stories from https://dev.to/
type DevToSource struct{}

// Fetch gets news from the Dev.To
func (l *DevToSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	resp, err := http.Get(DevToURL)
	if err != nil {
		return news, err
	}

	defer resp.Body.Close()

	doc := Rss{}
	d := xml.NewDecoder(resp.Body)

	if err := d.Decode(&doc); err != nil {
		return news, err
	}

	for i, item := range doc.Item {
		if i >= count {
			break
		}

		news[i] = item.Link
	}

	return news, nil
}

// DZoneSource fetches latest stories from http://feeds.dzone.com/home
type DZoneSource struct{}

// Fetch gets news from the DZone
func (l *DZoneSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	resp, err := http.Get(DZoneURL)
	if err != nil {
		return news, err
	}

	defer resp.Body.Close()

	doc := Rss{}
	d := xml.NewDecoder(resp.Body)

	if err := d.Decode(&doc); err != nil {
		return news, err
	}

	for i, item := range doc.Item {
		if i >= count {
			break
		}

		news[i] = item.Link
	}

	return news, nil
}

// LobstersSource fetches new stories from https://lobste.rs
type LobstersSource struct{}

// Fetch gets news from the Lobsters
func (l *LobstersSource) Fetch(count int) (map[int]string, error) {
	offset := float64(count) / float64(25)
	pages := int(math.Ceil(offset))
	news := make(map[int]string)
	newsIndex := 0

	for p := 1; p <= pages; p++ {
		url := fmt.Sprintf("%s/page/%d", LobstersURL, p)
		resp, err := http.Get(url)
		if err != nil {
			handleError(err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			handleError(err)
			continue
		}

		doc.Find(".link a.u-url").Each(func(_ int, s *goquery.Selection) {
			href, exist := s.Attr("href")
			if !exist {
				fmt.Println(red("can't find any stories..."))
			}

			if newsIndex >= count {
				return
			}

			// if internal link
			if strings.HasPrefix(href, "/") {
				href = LobstersURL + href
			}

			news[newsIndex] = href
			newsIndex++
		})

		resp.Body.Close()
	}

	return news, nil
}

// RedditSource fetches new stories from reddit.com/r/programming.
type RedditSource struct{}

// Fetch gets news from the Reddit
func (rs *RedditSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)

	s := geddit.NewSession(fmt.Sprintf("desktop:com.github.Bunchhieng.%s:%s", AppName, AppVersion))
	subs, err := s.SubredditSubmissions(
		"programming",
		geddit.HotSubmissions,
		geddit.ListingOptions{
			Count: count,
			Limit: count,
		},
	)

	if err != nil {
		return news, err
	}

	for i, sub := range subs {
		news[i] = sub.URL
	}

	return news, nil
}

// HackerNewsSource fetches new stories from news.ycombinator.com.
type HackerNewsSource struct{}

// Fetch gets news from the HackerNews
func (hn *HackerNewsSource) Fetch(count int) (map[int]string, error) {
	news := make(map[int]string)
	// 30 news per page
	pages := count / 30
	for i := 0; i <= pages; i++ {
		resp, err := http.Get(HackerNewsURL + strconv.Itoa(pages))
		if err != nil {
			handleError(err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			handleError(err)
			continue
		}

		doc.Find("a.storylink").Each(func(i int, s *goquery.Selection) {
			href, exist := s.Attr("href")
			if !exist {
				fmt.Println(red("can't find any stories..."))
			}
			news[i] = href
		})

		resp.Body.Close()
	}

	return news, nil
}
