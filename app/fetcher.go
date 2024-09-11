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
		title := strings.TrimSpace(s.Find(".link a").Text())
		url, _ := s.Find(".link a").Attr("href")
		url = strings.TrimSpace(url)
		votes := strings.TrimSpace(s.Find(".voters .score").Text())
		author := strings.TrimSpace(s.Find("a.u-author").Text())
		comments := strings.TrimSpace(s.Find(".comments_label a").Text())

		var tags []string
		s.Find(".tags a").Each(func(i int, s *goquery.Selection) {
			tags = append(tags, s.Text())
		})

		stories = append(stories, Story{
			Title:    title,
			Link:     url,
			Comments: comments,
			Author:   author,
			Votes:    votes,
			Tags:     tags,
		})
	})

	return stories, nil
}
