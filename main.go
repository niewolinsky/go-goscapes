package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sessionState is used to track which model is focused
type sessionState uint

func main() {
	p := tea.NewProgram(newModel(defaultTime))

	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota
	timerView2
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	state  sessionState
	timer  timer.Model
	timer2 timer.Model
}

func newModel(timeout time.Duration) mainModel {
	m := mainModel{state: timerView}
	m.timer = timer.New(timeout)
	m.timer2 = timer.New(timeout)
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.timer.Init(), m.timer2.Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == timerView {
				m.state = timerView2
			} else {
				m.state = timerView
			}
		case "n":
			if m.state == timerView {
				m.timer = timer.New(defaultTime)
				cmds = append(cmds, m.timer.Init())
			} else {
				m.timer2 = timer.New(defaultTime)
				cmds = append(cmds, m.timer2.Init())
			}
		}
		switch m.state {
		// update whichever model is focused
		case timerView2:
			m.timer2, cmd = m.timer2.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		}
	case timer.TickMsg:
		var cmd1 tea.Cmd
		var cmd2 tea.Cmd
		m.timer2, cmd1 = m.timer2.Update(msg)
		m.timer, cmd2 = m.timer.Update(msg)
		cmds = append(cmds, cmd1, cmd2)
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	if m.state == timerView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), modelStyle.Render(m.timer2.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), focusedModelStyle.Render(m.timer2.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == timerView {
		return "timer"
	}
	return "timer2"
}
