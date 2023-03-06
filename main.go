package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

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
	soundscape_02
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
	players []oto.Player
}

func newModel(timeout time.Duration, players []oto.Player) mainModel {
	m := mainModel{state: soundscape_01}
	m.players = players
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == soundscape_01 {
				m.state = soundscape_02
			} else {
				m.state = soundscape_01
			}
		case "enter":
			if m.state == soundscape_01 {
				if m.players[0].IsPlaying() {
					m.players[0].Pause()
					m.players[0].(io.Seeker).Seek(0, io.SeekStart)
				} else {
					m.players[0].Play()
				}
			} else if m.state == soundscape_02 {
				if m.players[1].IsPlaying() {
					m.players[1].Pause()
					m.players[1].(io.Seeker).Seek(0, io.SeekStart)
				} else {
					m.players[1].Play()
				}
			}
		case "i":
			if m.state == soundscape_01 {
				if m.players[0].Volume() > 0.1 {
					m.players[0].SetVolume(m.players[0].Volume() - 0.1)
				}
			} else if m.state == soundscape_02 {
				if m.players[1].Volume() > 0.1 {
					m.players[1].SetVolume(m.players[1].Volume() - 0.1)
				}
			}

		case "k":
			if m.state == soundscape_01 {
				if m.players[0].Volume() < 1.0 {
					m.players[0].SetVolume(m.players[0].Volume() + 0.1)
				}
			} else if m.state == soundscape_02 {
				if m.players[1].Volume() < 1.0 {
					m.players[1].SetVolume(m.players[1].Volume() + 0.1)
				}
			}
		}
	}
	return m, nil
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	if m.state == soundscape_01 {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render("rain"), modelStyle.Render("bells"))
	} else if m.state == soundscape_02 {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("rain"), focusedModelStyle.Render("bells"))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: play %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == soundscape_01 {
		return "soundscape_01"
	} else if m.state == soundscape_02 {
		return "soundscape_02"
	}

	return ""
}
