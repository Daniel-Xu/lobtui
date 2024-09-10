package app

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchStories(app *App) ([]Story, error) {
	url := fmt.Sprintf("https://lobste.rs/page/%d", app.page)
	resp, err := app.client.R().Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return nil, err
	}

	var stories []Story
	doc.Find(".story").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".link a").Text()
		url, _ := s.Find(".link a").Attr("href")
		stories = append(stories, Story{Title: title, Link: url})
	})

	return stories, nil
}
