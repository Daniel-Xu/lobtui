package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// key handling for list item
type delegateKeyMap struct {
	open key.Binding
	show key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		open: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Open URL"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.open):
				return m.NewStatusMessage("You chose " + title)
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

// ey handling for list

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
