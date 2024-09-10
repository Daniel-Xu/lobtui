package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-resty/resty/v2"
)

const footerHeight = 3

type VerticalLayout struct {
	body   viewport.Model
	footer string
}

type App struct {
	layout  VerticalLayout
	stories []Story
	cursor  int
	page    int
	client  *resty.Client
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
		layout: VerticalLayout{
			body:   viewport.New(100, 100), // Initialize with a large size, it will be adjusted later
			footer: "Page: 1 | Press 'n' for next, 'b' for previous, 'q' to quit",
		},
		stories: []Story{},
		cursor:  0,
		page:    1,
		client:  resty.New(),
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	case storiesMsg:
		a.stories = msg
		a.updateContent()
	case tea.WindowSizeMsg:
		a.layout.body.Width = msg.Width
		a.layout.body.Height = msg.Height - footerHeight
		a.updateContent()
		return a, nil // Return immediately after updating window size
	}
	a.layout.body, cmd = a.layout.body.Update(msg)
	return a, cmd
}

func (a *App) View() string {
	return fmt.Sprintf("%s\n%s", a.layout.body.View(), a.layout.footer)
}

func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return a, tea.Quit
	case "up", "k":
		if a.cursor > 0 {
			a.cursor--
			a.updateContent()
		}
	case "down", "j":
		if a.cursor < len(a.stories)-1 {
			a.cursor++
			a.updateContent()
		}
	case "n":
		a.page++
		return a, a.fetchNextPage
	case "b", "p":
		if a.page > 1 {
			a.page--
			return a, a.fetchPreviousPage
		}
	}
	return a, nil
}

func (a *App) updateContent() {
	var content strings.Builder
	for i, story := range a.stories {
		item := fmt.Sprintf("%d. %s", i+1, story.Title)
		if a.cursor == i {
			item = selectedItemStyle.Render("> " + item)
		} else {
			item = regularItemStyle.Render("  " + item)
		}
		content.WriteString(item + "\n")
	}

	renderedContent := listStyle.Render(content.String())
	a.layout.body.SetContent(renderedContent)
	a.layout.footer = fmt.Sprintf("Page: %d | Press 'n' for next, 'b' for previous, 'q' to quit", a.page)
}

func (a *App) fetchNextPage() tea.Msg {
	stories, err := fetchStories(a)
	if err != nil {
		a.page--
		return nil
	}
	return storiesMsg(stories)
}

func (a *App) fetchPreviousPage() tea.Msg {
	stories, err := fetchStories(a)
	if err != nil {
		a.page++
		return nil
	}
	return storiesMsg(stories)
}

var (
	listStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, false).
			MarginLeft(2).Align(lipgloss.Left).Width(100)

	selectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#AA55AA40")).
				Bold(true)

	regularItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))
)
