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

type customPlayer struct {
	player   oto.Player
	isActive bool
}

func main() {
	otoCtx, readyChan, err := oto.NewContext(44100, 2, 2)
	if err != nil {
		log.Fatal(err)
	}
	<-readyChan

	players := map[int]customPlayer{}

	soundscapes, err := Files.ReadDir("soundscapes")
	if err != nil {
		log.Fatal(err)
	}
	for i, soundscape := range soundscapes {
		fileBytes, err := Files.ReadFile(fmt.Sprintf("soundscapes/%s", soundscape.Name()))
		if err != nil {
			log.Fatal(err)
		}

		fileBytesReader := bytes.NewReader(fileBytes)

		decodedMp3, err := mp3.NewDecoder(fileBytesReader)
		if err != nil {
			log.Fatal(err)
		}

		player := customPlayer{otoCtx.NewPlayer(decodedMp3), false}
		players[i] = player
	}

	for i := range players {
		defer players[i].player.Close()
	}

	p := tea.NewProgram(newModel(players), tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const (
	soundscape_01 sessionState = iota
	soundscape_02
	soundscape_03
	soundscape_04
	soundscape_05
	soundscape_06
	soundscape_07
	soundscape_08
	soundscape_09
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
	players map[int]customPlayer
}

func newModel(players map[int]customPlayer) mainModel {
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
			if m.state == 0 {
				m.state = soundscape_02
			} else if m.state == 1 {
				m.state = soundscape_03
			} else if m.state == 2 {
				m.state = soundscape_04
			} else if m.state == 3 {
				m.state = soundscape_05
			} else if m.state == 4 {
				m.state = soundscape_06
			} else if m.state == 5 {
				m.state = soundscape_07
			} else if m.state == 6 {
				m.state = soundscape_08
			} else if m.state == 7 {
				m.state = soundscape_09
			} else if m.state == 8 {
				m.state = soundscape_01
			}
		case "enter":
			if m.players[int(m.state)].player.IsPlaying() {
				m.players[int(m.state)].player.Pause()
				m.players[int(m.state)].player.(io.Seeker).Seek(0, io.SeekStart)

				if entry, ok := m.players[int(m.state)]; ok {
					entry.isActive = false
					m.players[int(m.state)] = entry
				}
			} else {
				m.players[int(m.state)].player.Play()

				if entry, ok := m.players[int(m.state)]; ok {
					entry.isActive = true
					m.players[int(m.state)] = entry
				}

				return m, keepAlive(m.players[int(m.state)].player)
			}
		case "i":
			if m.players[int(m.state)].player.Volume() > 0.1 {
				fmt.Println(m.players[int(m.state)].player.Volume())
				m.players[int(m.state)].player.SetVolume(m.players[int(m.state)].player.Volume() - 0.1)
				fmt.Println(m.players[int(m.state)].player.Volume())
			}
		case "k":
			if m.players[int(m.state)].player.Volume() < 1.0 {
				fmt.Println(m.players[int(m.state)].player.Volume())
				m.players[int(m.state)].player.SetVolume(m.players[int(m.state)].player.Volume() + 0.1)
				fmt.Println(m.players[int(m.state)].player.Volume())
			}
		case "p":
			return m, keepAlive(m.players[0].player)
		}
	case statusMsg:
		if !m.players[int(m.state)].isActive {
			return m, nil
		}
		m.players[int(m.state)].player.(io.Seeker).Seek(0, io.SeekStart)
		m.players[int(m.state)].player.Play()
		return m, keepAlive(m.players[int(m.state)].player)
	}
	return m, nil
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()

	switch m.state {
	case soundscape_01:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render("RAIN"), modelStyle.Render("THUNDER"), modelStyle.Render("WAVES")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("WIND"), modelStyle.Render("FIRE"), modelStyle.Render("BIRDS")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("CRICKETS"), modelStyle.Render("BOWLS"), modelStyle.Render("WHITE NOISE")))
	case soundscape_02:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), focusedModelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_03:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), focusedModelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_04:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_05:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), focusedModelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_06:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), focusedModelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_07:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_08:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), focusedModelStyle.Render("🔔"), modelStyle.Render("🐦")))
	case soundscape_09:
		s += lipgloss.JoinVertical(lipgloss.Top, lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), modelStyle.Render("🐦")), lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("🌧️"), modelStyle.Render("🔔"), focusedModelStyle.Render("🐦")))
	}

	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: play %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	//?smarter way to do that?
	switch m.state {
	case soundscape_01:
		return "rain"
	case soundscape_02:
		return "thunder"
	case soundscape_03:
		return "waves"
	case soundscape_04:
		return "wind"
	case soundscape_05:
		return "fire"
	case soundscape_06:
		return "birds"
	case soundscape_07:
		return "crickets"
	case soundscape_08:
		return "bowls"
	case soundscape_09:
		return "noise"
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
