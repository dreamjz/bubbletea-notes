package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// messages that implement the interface tea.Msg
// tes.Msg is an interface with no methods
// no need to implement any method
type (
	gotRepoSuccessMsg []repo
	gotRepoErrMsg     error
)

type repo struct {
	Name string `json:"name"`
}

const reposURL = "https://api.github.com/orgs/charmbracelet/repos"

// bubbletea command func tea.Cmd
func getRepos() tea.Msg {
	// create new http get request
	req, err := http.NewRequest(http.MethodGet, reposURL, nil)
	if err != nil {
		return gotRepoErrMsg(err) // wrap err
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	// send get request with default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gotRepoErrMsg(err)
	}
	defer resp.Body.Close() // close resource

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return gotRepoErrMsg(err)
	}

	var repos []repo

	err = json.Unmarshal(data, &repos)
	if err != nil {
		return gotRepoErrMsg(err)
	}

	return gotRepoSuccessMsg(repos) // wrap []repo
}

// bubbletea model

type model struct {
	textInput textinput.Model // bubbles text input component, akin to <input type="text"> in HTML
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = "charmbracelet/"
	ti.Placeholder = "repo..."
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	// ti.ShowSuggestions = true // there is no this field in bubbles v0.16.1
	return model{textInput: ti}
}

func (m model) Init() tea.Cmd {
	// use tea.Batch to return multiple commands
	return tea.Batch(getRepos, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case gotRepoSuccessMsg:
		// m.textInput.SetValue("success")
		// auto complete not available in bubbles v0.16.1
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"What's your favorite Charm repository?\n\n%s\n\n%\n",
		m.textInput.View(),
		"(tab to complete, ctrl+n/ctrl+p to cycle through suggestion, esc to quit",
	)
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		log.Fatal(err)
	}
}
