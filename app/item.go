package app

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Story struct {
	Title        string
	Link         string
	Author       string
	Tags         []string
	Votes        int
	CommentCount string
	CommentLink  string
}

type storiesMsg []Story
