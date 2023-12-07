package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	// get color function from env color profile to create color from string
	color = termenv.EnvColorProfile().Color
	// functions to style a string
	keyword = termenv.Style{}.Foreground(color("204")).Background(color("235")).Styled
	help    = termenv.Style{}.Foreground(color("241")).Styled
)

type model struct {
	altscreen bool // enable alternative screen
	quiting   bool // quit program
}

// implement interface tea.Model

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quiting = true
			return m, tea.Quit
		case " ":
			var cmd tea.Cmd
			if m.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altscreen = !m.altscreen
			return m, cmd
		}
	}

	return m, nil
}

// View is used to render UI
func (m model) View() string {
	if m.quiting {
		return "Bye!\n"
	}

	const (
		altscreenMode = " altscreen mode "
		inlineMode    = " inline mode "
	)

	var mode string
	if m.altscreen {
		mode = altscreenMode
	} else {
		mode = inlineMode
	}

	return fmt.Sprintf(" You're in %s", keyword(mode)) +
		help(" space: witch modes / q: exit\n")
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
