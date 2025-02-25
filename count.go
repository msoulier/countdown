package main

// A simple example that shows how to render a progress bar in a "pure"
// fashion. In this example we bump the progress by 25% every second,
// maintaining the progress state on our top level model using the progress bar
// model's ViewAs method only for rendering.
//
// The signature for ViewAs is:
//
//     func (m Model) ViewAs(percent float64) string
//
// So it takes a float between 0 and 1, and renders the progress bar
// accordingly. When using the progress bar in this "pure" fashion and there's
// no need to call an Update method.
//
// The progress bar is also able to animate itself, however. For details see
// the progress-animated example.

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

const (
	Countdown string = "countdown"
	Countup          = "countup"
)

var (
	endseconds  int64
	curseconds  int64
	progressbar bool
	fromcolour  = "#0000FF"
	tocolour    = "#FF0000"
	mode        = "countup"
)

func init() {
	flag.Int64Var(&endseconds, "seconds", 0, "Number of seconds to count down")
	flag.BoolVar(&progressbar, "prog", true, "Display an in-colour progress bar")
	flag.StringVar(&fromcolour, "fromcolour", "#0000FF", "Left-hand colour of gradient")
	flag.StringVar(&tocolour, "tocolour", "#FF0000", "Right-hand colour of gradient")
	flag.StringVar(&mode, "mode", "countup", "Should the progress bar count up or down")
	flag.Parse()
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

type model struct {
	percent    float64
	endseconds int64
	curseconds int64
	progress   progress.Model
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		m.curseconds += 1
		if m.curseconds >= m.endseconds {
			m.percent = 1.0
			return m, tea.Quit
		}
		m.percent = float64(m.curseconds) / float64(m.endseconds)

		return m, tickCmd()

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func main() {
	if progressbar {
		prog := progress.New(progress.WithScaledGradient(fromcolour, tocolour))
		mod := model{
			progress:   prog,
			endseconds: endseconds,
			curseconds: 0,
		}

		if _, err := tea.NewProgram(mod).Run(); err != nil {
			fmt.Println("Oh no!", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Here I will put a plain old ASCII countdown")
		os.Exit(1)
	}
}
