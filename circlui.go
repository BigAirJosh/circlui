package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bigairjosh/circlui/api"
	table "github.com/calyptia/go-bubble-table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jszwedko/go-circleci"
	"golang.org/x/term"
)


var cliConfig api.CliConfig = api.LoadCliConfig()
var client circleci.Client = circleci.Client{Token: cliConfig.Token}

func main() {

  // builds, _ := client.ListRecentBuilds(-1, 0)

  // for _, b := range builds {
  //   fmt.Printf("%s\n", b.Reponame)
  // }

	err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	styleDoc = lipgloss.NewStyle().
    Border(lipgloss.NormalBorder(), true).
    Padding(1)
)

func initialModel() model {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		w = 80
		h = 24
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom

	tbl := table.New([]string{"ID", "App", "Branch", "Remaining"}, w, h)
	tbl.Styles = table.Styles{
		Title:       lipgloss.NewStyle().Bold(true),
		SelectedRow: lipgloss.NewStyle().Foreground(lipgloss.Color("#AFD75F")),
	}
  
  builds, _ := client.ListRecentBuilds(100, 0)
	rows := buildsToRows(builds)
 	tbl.SetRows(rows)

  model := model{table: tbl, builds: builds}

  return model
}

func buildsToRows(builds []*circleci.Build) []table.Row {
	rows := make([]table.Row, len(builds))
	for i := 0; i < len(builds); i++ {
		rows[i] = table.SimpleRow{
			builds[i].BuildNum,
			builds[i].JobName,
			builds[i].Branch,
			builds[i].BuildTimeMillis,
		}
	}
  return rows
}

type model struct {
	table table.Model
  builds []*circleci.Build
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
  return tea.Every(time.Second,
  func(t time.Time) tea.Msg {
    return TickMsg(t)
  })
}

func (m model) Init() tea.Cmd {
	return tickEvery()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := styleDoc.GetPadding()
		m.table.SetSize(
			msg.Width-left-right,
			msg.Height-top-bottom,
		)
  case TickMsg:
    // for i, j := range m.jobs {
    //   if j.RemianingSeconds > 0 {
    //     m.jobs[i].RemianingSeconds = j.RemianingSeconds - 1
    //   }
    // }
    // m.table.SetRows(jobsToRows(m.jobs))
    return m, tickEvery()
	case tea.KeyMsg:
		switch msg.String() {
    case "ctrl+r":
    m.table.SetRows(buildsToRows(m.builds))
    return m, nil
    case "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(
		m.table.View(),
	)
}
