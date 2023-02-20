package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bigairjosh/circlui/api"
	table "github.com/calyptia/go-bubble-table"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jszwedko/go-circleci"
	"golang.org/x/term"
  "github.com/hako/durafmt"
)

var cliConfig api.CliConfig = api.LoadCliConfig()
var client circleci.Client = circleci.Client{Token: cliConfig.Token}

type keyMap struct {
	Reload key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Reload},       // first column
		{k.Help}, // second column
    {k.Quit},
	}
}

var keys = keyMap{
	Reload: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "reload"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

var units, _ = durafmt.DefaultUnitsCoder.Decode("y:y,w:w,d:d,h:h,m:m,s:s,ms:ms,mi:mi")

func main() {

  client.Debug = false
	err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getBranch() string {
  path, _ := os.Getwd()
  cmd := exec.Command("git", "branch", "--show-current")
  cmd.Dir = path
  out, _ := cmd.Output()

  return strings.TrimSpace(string(out))
}

func getRepository() string {
  path, _ := os.Getwd()
  cmd := exec.Command("git", "rev-parse", "--show-toplevel")
  cmd.Dir = path
  out, _ := cmd.Output()

  splitPath := strings.Split(string(out), "/")
  repo := splitPath[len(splitPath) - 1]

  return strings.TrimSpace(repo)
}

var (
	styleDoc = lipgloss.NewStyle().
		// Border(lipgloss.NormalBorder(), true).
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
	h = h - top - bottom - 1

	tbl := table.New([]string{"Pipeline", "Workflow", "Status", "Branch", "Commiter", "Start", "Duration"}, w, h - 6)
	tbl.Styles = table.Styles{
		Title:       lipgloss.NewStyle().Bold(true),
		SelectedRow: lipgloss.NewStyle().Foreground(lipgloss.Color("#AFD75F")),
	}

	model := model{
		table:      tbl,
		builds:     nil,
    branchFilter: getBranch(),
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
    repository: getRepository(),
	}

  loadBuilds(&model)

	return model
}

func loadBuilds(model *model) {
  builds, _ := client.ListRecentBuildsForProject("deliveroo", model.repository, model.branchFilter, "", 100, 0)
  model.builds = builds
	model.table.SetRows(buildsToRows(builds))
}

func buildsToRows(builds []*circleci.Build) []table.Row {
	rows := make([]table.Row, len(builds))
	for i := 0; i < len(builds); i++ {
		startTime := time.Since(time.Now()).Round(1 * time.Minute)
		if builds[i].StartTime != nil {
			startTime = time.Since(*builds[i].StartTime).Round(1 * time.Minute)
		}
    var duration time.Duration 
    if builds[i].Status == "running" {
      duration = time.Since(*builds[i].StartTime)
    }
		if builds[i].StopTime != nil {
      duration = builds[i].StopTime.Sub(*builds[i].StartTime)
		}

    comitter := ""
    if(len(builds[i].AllCommitDetails) > 0) {
      comitter = builds[i].AllCommitDetails[0].CommitterLogin
    }

		rows[i] = table.SimpleRow{
			builds[i].Reponame,
			builds[i].Workflows.JobName,
			builds[i].Status,
			builds[i].Branch,
      comitter,
			durafmt.ParseShort(startTime).String(),
			durafmt.Parse(duration).LimitFirstN(2).Format(units),
		}
	}
	return rows
}

type model struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	table      table.Model
	builds     []*circleci.Build
  branchFilter string
  repository string
}

type TickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Every(time.Second,
		func(t time.Time) tea.Msg {
			return TickMsg(t)
		})
}

type RefreshMsg time.Time

func refreshEvery() tea.Cmd {
	return tea.Every(time.Second * 30,
		func(t time.Time) tea.Msg {
			return RefreshMsg(t)
		})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickEvery(), refreshEvery())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := styleDoc.GetPadding()
		m.table.SetSize(
			msg.Width-left-right,
			msg.Height-top-bottom-1,// TODO break out to function
		)
		m.help.Width = msg.Width
	case TickMsg:
    m.table.SetRows(buildsToRows(m.builds))
		return m, tickEvery()
	case RefreshMsg:
    loadBuilds(&m)
		return m, refreshEvery()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Reload):
      loadBuilds(&m)
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
      return m, nil
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
    case msg.String() == "backspace":
      l := len(m.branchFilter)
      if(l > 0) {
        m.branchFilter = m.branchFilter[:l-1]
      }
      return m, nil
    case msg.String() == "enter":
      loadBuilds(&m)
      return m, nil
    default:
      m.branchFilter += msg.String()
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	helpView := m.help.View(m.keys)
	height := strings.Count(helpView, "\n")

	return styleDoc.Render(m.table.View()) + "\nbranch filter: " + m.branchFilter + "\n" + strings.Repeat("\n", height) + helpView
}
