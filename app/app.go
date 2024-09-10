package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-resty/resty/v2"
)

type App struct {
	viewport viewport.Model
	stories  []Story
	cursor   int
	page     int
	client   *resty.Client
}

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

func NewApp() *App {
	return &App{
		viewport: viewport.New(80, 20),
		stories:  []Story{},
		cursor:   0,
		page:     1,
		client:   resty.New(),
	}
}

func (a *App) Init() tea.Cmd {
	stories, err := fetchStories(a)
	if err != nil {
		return tea.Quit
	}
	a.stories = stories
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	case storiesMsg:
		a.stories = msg
		return a, nil
	case tea.WindowSizeMsg:
		a.viewport.Width = msg.Width
		a.viewport.Height = msg.Height - 1
		return a, nil
	}
	return a, nil
}

func (a *App) View() string {
	var content string
	for i, story := range a.stories {
		cursor := " "
		if a.cursor == i {
			cursor = ">"
		}
		content += fmt.Sprintf("%s %s\n", cursor, story.Title)
	}
	a.viewport.SetContent(content)
	return a.viewport.View()
}

func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return a, tea.Quit
	case "up", "k":
		if a.cursor > 0 {
			a.cursor--
		}
	case "down", "j":
		if a.cursor < len(a.stories)-1 {
			a.cursor++
		}
	case "n":
		a.page++
		stories, err := fetchStories(a)
		if err != nil {
			a.page--
			return a, nil
		}
		a.stories = stories
		return a, nil
	case "b", "p":
		if a.page > 1 {
			a.page--
			stories, err := fetchStories(a)
			if err != nil {
				a.page++
				return a, nil
			}
			a.stories = stories
			return a, nil
		}
	}
	return a, nil
}
