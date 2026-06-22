package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"example.com/m/v2/pkg"
)

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(false)

type KeyMap struct {
	table.KeyMap
	Kill   key.Binding
	Cancel key.Binding
}

func (km KeyMap) ShortHelp() []key.Binding {
	short := km.KeyMap.ShortHelp()
	return append(short, km.Kill, km.Cancel)
}

func (km KeyMap) FullHelp() [][]key.Binding {
	full := km.KeyMap.FullHelp()
	return append(full, []key.Binding{km.Kill, km.Cancel})
}

func KeyMaps() KeyMap {
	return KeyMap{
		KeyMap: table.DefaultKeyMap(),
		Kill: key.NewBinding(
			key.WithKeys("enter", "x"),
			key.WithHelp("x/enter", "Kill process"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc", "q"),
			key.WithHelp("esc/q", "quit"),
		),
	}
}

type model struct {
	table     table.Model
	processes []pkg.RunningProcess
	keys      KeyMap
	help      help.Model
}

func (m model) Init() tea.Cmd { return nil }

type processRefresMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Cancel):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Kill):
			return m, func() tea.Msg {
				process := m.processes[m.table.Cursor()]
				process.KillProcess()
				return processRefresMsg{}
			}
		}
	case processRefresMsg:
		processes, err := pkg.GetRunningProcesses()
		if err != nil {
			return m, tea.Quit
		}
		m.processes = []pkg.RunningProcess(processes)
		m.table.SetRows(generateRows(m.processes))
	}
	m.table, cmd = m.table.Update(msg)
	return m, tea.Batch(cmd, pollProcesses())
}

func (m model) View() tea.View {
	return tea.NewView(
		baseStyle.Render(m.table.View()) + "\n" + m.help.View(m.keys),
	)
}

func generateRows(processes []pkg.RunningProcess) []table.Row {
	rows := []table.Row{}
	for _, val := range processes {
		rows = append(rows, table.Row{
			val.COMMAND,
			val.PID,
			val.USER,
			val.FD,
			val.TYPE,
			strconv.Itoa(val.DEVICE),
			val.SIZE,
			val.NODE,
			strconv.Itoa(val.PORT),
			val.NAME,
		})
	}
	return rows
}
func pollProcesses() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return processRefresMsg{}
	})
}

func main() {
	processes, err := pkg.GetRunningProcesses()
	if err != nil {
		fmt.Println(err)
		return
	}
	columns := []table.Column{
		{Title: "Command", Width: 10},
		{Title: "PID", Width: 10},
		{Title: "User", Width: 10},
		{Title: "FD", Width: 10},
		{Title: "TYPE", Width: 10},
		{Title: "Device", Width: 10},
		{Title: "Size", Width: 10},
		{Title: "Node", Width: 10},
		{Title: "Port", Width: 15},
		{Title: "Name", Width: 10},
	}
	rows := generateRows(processes)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(120),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.BorderBottom(false)
	t.SetStyles(s)

	m := model{
		table: t,
		processes: processes,
		keys: KeyMaps(),
		help: help.New(),
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
