package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-resty/resty/v2"
)

// application state
type App struct {
	stories []Story
	page    int
	client  *resty.Client
	list    list.Model
	keys    *listKeyMap
}

func NewApp() *App {
	delegateKeys := newDelegateKeyMap()
	delegate := newItemDelegate(delegateKeys)

	storyList := list.New([]list.Item{}, delegate, 0, 0)
	storyList.SetShowStatusBar(false)
	storyList.SetFilteringEnabled(false)
	listKeys := newListKeyMap()

	storyList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.nextPage,
			listKeys.prevPage,
		}
	}

	return &App{
		page:   1,
		client: resty.New(),
		list:   storyList,
		keys:   listKeys,
	}
}

func (a *App) Init() tea.Cmd {
	stories, err := fetchStories(a)
	a.setItems(stories)

	if err != nil {
		return tea.Quit
	}
	return nil
}

func (a *App) setItems(stories []Story) {
	items := make([]list.Item, len(stories))
	for i, story := range stories {
		items[i] = item{
			title: fmt.Sprintf("%-2d %s", i+1, story.Title),
			desc:  fmt.Sprintf("Author: %s | Votes: %d | Comments: %s", story.Author, story.Votes, story.CommentCount),
		}
	}
	a.list.SetItems(items)
	a.list.ResetSelected()
	a.list.Title = fmt.Sprintf("Lobsters Stories - Page %d", a.page)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.nextPage):
			a.page++
			return a, a.fetchNextPage

		case key.Matches(msg, a.keys.prevPage):
			if a.page > 1 {
				a.page--
				return a, a.fetchPreviousPage
			}
		}

	case storiesMsg:
		a.setItems(msg)

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		a.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newListModel, cmd := a.list.Update(msg)
	a.list = newListModel
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	return docStyle.Render(a.list.View())
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
	docStyle = lipgloss.NewStyle().Margin(1, 2)
)
