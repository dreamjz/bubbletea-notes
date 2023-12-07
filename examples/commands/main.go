package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

const url = "https://charm.sh/"

type model struct {
	spinner spinner.Model
	status  int
	err     error
}

type statusMsg int

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	}

	return statusMsg(res.StatusCode)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(checkServer, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case statusMsg:
		m.status = int(msg)
		return m, tea.Quit
	case errMsg:
		m.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Checking %s ...%s\n", url, m.spinner.View())

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n"
}

func main() {
	m := model{
		spinner: spinner.New(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
