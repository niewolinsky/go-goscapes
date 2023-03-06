package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

// sessionState is used to track which model is focused
type sessionState uint

type statusMsg uint

//go:embed "soundscapes"
var Files embed.FS

func main() {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		log.Fatal(err)
	}
	<-readyChan

	players := []oto.Player{}

	soundscapes, err := Files.ReadDir("soundscapes")
	if err != nil {
		log.Fatal(err)
	}
	for _, soundscape := range soundscapes {
		fileBytes, err := Files.ReadFile(fmt.Sprintf("soundscapes/%s", soundscape.Name()))
		if err != nil {
			log.Fatal(err)
		}

		fileBytesReader := bytes.NewReader(fileBytes)

		decodedMp3, err := mp3.NewDecoder(fileBytesReader)
		if err != nil {
			log.Fatal(err)
		}

		player := otoCtx.NewPlayer(decodedMp3)
		players = append(players, player)
	}

	for i := range players {
		defer players[i].Close()
	}

	p := tea.NewProgram(newModel(players))
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
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("241"))
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

func newModel(players []oto.Player) mainModel {
	m := mainModel{state: soundscape_01}
	m.players = players
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == 1 {
				m.state = soundscape_02
			} else {
				m.state = soundscape_01
			}
		case "enter":
			if m.players[m.state-1].IsPlaying() {
				m.players[m.state-1].Pause()
				m.players[m.state-1].(io.Seeker).Seek(0, io.SeekStart)
			} else {
				m.players[m.state-1].Play()
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

		case "p":
			return m, keepAlive(m.players[0])
		}
	case statusMsg:
		fmt.Println("wykonuje sie")
		m.players[0].(io.Seeker).Seek(0, io.SeekStart)
		m.players[0].Play()
		return m, keepAlive(m.players[0])
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

func keepAlive(player oto.Player) tea.Cmd {
	return func() tea.Msg {
		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
		}

		return statusMsg(0)
	}
}
