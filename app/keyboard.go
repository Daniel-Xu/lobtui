package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
)

// key handling for list item
type delegateKeyMap struct {
	open key.Binding
	show key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		open: key.NewBinding(
			key.WithKeys("enter", "o"),
			key.WithHelp("enter/o", "Open URL"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var url string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
			url = i.URL()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.open):
				browser.Stdout = nil
				browser.Stderr = nil
				if err := browser.OpenURL(url); err != nil {
					return m.NewStatusMessage("Failed to open URL: " + err.Error())
				}
				return m.NewStatusMessage("Opening " + title)
			}
		}

		return nil
	}

	help := []key.Binding{keys.open}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	return d
}

// key handling for list
type listKeyMap struct {
	nextPage key.Binding
	prevPage key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		nextPage: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next remote page"),
		),
		prevPage: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "previous remote page"),
		),
	}
}
