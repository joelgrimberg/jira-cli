package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	choice string
}

// Define a struct to match the JSON structure
type JiraResponse struct {
	Issues []struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
		} `json:"fields"`
	} `json:"issues"`
}

// Function to load the file and export summaries
func LoadSummariesFromFile(filename string) ([]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the file content
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into JiraResponse struct
	var jiraResponse JiraResponse
	err = json.Unmarshal(byteValue, &jiraResponse)
	if err != nil {
		return nil, err
	}

	// Collect summaries
	var summaries []string
	for _, issue := range jiraResponse.Issues {
		summaries = append(summaries, issue.Fields.Summary)
	}

	return summaries, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return fmt.Sprintf("You've selected: %s", m.choice)
	}

	return docStyle.Render(m.list.View())
}

func main() {
	summaries, err := LoadSummariesFromFile("cmd/jira-cli/testdata/jira-example.json")
	if err != nil {
		fmt.Println("Error loading summaries:", err)
		return
	}

	// Initialize an empty slice of list.Item
	items := []list.Item{}

	for _, summary := range summaries {
		// Create a new list.Item object with Title and Desc fields
		// Ensure Title and Desc are exported fields in the list.Item struct
		item := item{title: summary, desc: ""}
		items = append(items, item)
	}
	fmt.Printf("%v", items)

	config, err := LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	if config.RUN_MODE == "test" {
		fmt.Printf("Runmode: %s\n", config.RUN_MODE)
	}

	// p := tea.NewProgram(initialModel(summaries))
	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Tickets"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas there's been an error: %v\n", err)
		os.Exit(1)
	}
}
