package main

import (
	"fmt"

	"github.com/Daniel-Xu/lobtui/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(app.NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
