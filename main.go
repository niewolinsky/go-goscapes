package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

// sessionState is used to track which model is focused
type sessionState uint

func main() {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		log.Fatal(err)
	}
	<-readyChan

	mp3_01, err := os.Open("./01_rain.mp3")
	if err != nil {
		panic("opening my-file.mp3 failed: " + err.Error())
	}
	defer mp3_01.Close()

	mp3_02, err := os.Open("./02_bells.mp3")
	if err != nil {
		panic("opening my-file.mp3 failed: " + err.Error())
	}
	defer mp3_02.Close()

	mp3_01_dec, err := mp3.NewDecoder(mp3_01)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}
	mp3_02_dec, err := mp3.NewDecoder(mp3_02)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	play_01 := otoCtx.NewPlayer(mp3_01_dec)
	defer play_01.Close()
	play_02 := otoCtx.NewPlayer(mp3_02_dec)
	defer play_02.Close()

	players := []oto.Player{play_01, play_02}

	p := tea.NewProgram(newModel(defaultTime, players))
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const (
	defaultTime                = time.Minute
	soundscape_01 sessionState = iota
	sounsdcape_02
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
	state   sessionState
	timer   timer.Model
	timer2  timer.Model
	players []oto.Player
}

func newModel(timeout time.Duration, players []oto.Player) mainModel {
	m := mainModel{state: soundscape_01}
	m.timer = timer.New(timeout)
	m.timer2 = timer.New(timeout)
	m.players = players
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
			if m.state == soundscape_01 {
				m.state = sounsdcape_02
			} else {
				m.state = soundscape_01
			}
		case "n":
			if m.state == soundscape_01 {
				m.timer = timer.New(defaultTime)

				if m.players[0].IsPlaying() {
					m.players[0].Pause()
				} else {
					m.players[0].Play()
				}

				cmds = append(cmds, m.timer.Init())
			} else {
				m.timer2 = timer.New(defaultTime)

				if m.players[1].IsPlaying() {
					m.players[1].Pause()
				} else {
					m.players[1].Play()
				}

				cmds = append(cmds, m.timer2.Init())
			}
		}
		switch m.state {
		// update whichever model is focused
		case sounsdcape_02:
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
	if m.state == soundscape_01 {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), modelStyle.Render(m.timer2.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), focusedModelStyle.Render(m.timer2.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == soundscape_01 {
		return "timer"
	}

	return "timer2"
}
