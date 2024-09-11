package app

type item struct {
	title, desc, url, votes string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func (i item) URL() string         { return i.url }

type Story struct {
	Title    string
	Link     string
	Author   string
	Tags     []string
	Votes    string
	Comments string
}

type storiesMsg []Story
