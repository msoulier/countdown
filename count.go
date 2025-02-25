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
    mlib "github.com/msoulier/mlib"
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
	total_seconds  int64
	curseconds  int64
	progressbar bool
	fromcolour  = "#0000FF"
	tocolour    = "#FF0000"
	mode        = "countup"
)

func init() {
	flag.Int64Var(&total_seconds, "seconds", 0, "Number of seconds to count down")
	flag.BoolVar(&progressbar, "prog", true, "Display an in-colour progress bar")
	flag.StringVar(&fromcolour, "fromcolour", "#0000FF", "Left-hand colour of gradient")
	flag.StringVar(&tocolour, "tocolour", "#FF0000", "Right-hand colour of gradient")
	flag.StringVar(&mode, "mode", "countup", "Should the progress bar count up or down")
	flag.Parse()

    if total_seconds == 0 {
        flag.PrintDefaults()
        os.Exit(1)
    }

    // mode should be Countdown or Countup
    if mode != Countdown && mode != Countup {
        flag.PrintDefaults()
        os.Exit(1)
    }

    if mode == Countup {
        curseconds = 0
    } else if mode == Countdown {
        curseconds = total_seconds
    } else {
        panic("Unknown mode")
    }
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time

type model struct {
	percent    float64
	total_seconds int64
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
        if mode == Countup {
            m.curseconds += 1
            if m.curseconds >= m.total_seconds {
                m.percent = 1.0
                return m, tea.Quit
            }
            m.percent = float64(m.curseconds) / float64(m.total_seconds)
        } else if mode == Countdown {
            m.curseconds -= 1
            if m.curseconds == 0 {
                m.percent = 0.0
                return m, tea.Quit
            }
            m.percent = float64(m.curseconds) / float64(m.total_seconds)
        } else {
            panic("Unknown mode")
        }

		return m, tickCmd()

	default:
		return m, nil
	}
}

func (m model) View() string {
    current_duration := time.Duration(int64(time.Second) * m.curseconds)
    total_duration := time.Duration(int64(time.Second) * m.total_seconds)
    human_current := mlib.Duration2Human(current_duration, true)
    human_total := mlib.Duration2Human(total_duration, true)
    caption := ""
    if mode == Countdown {
        caption = fmt.Sprintf("Time remaining: %s", human_current)
    } else if mode == Countup {
        caption = fmt.Sprintf("Current time: %s\nUntil: %s", human_current, human_total)
    } else {
        panic("Unknown mode")
    }
	pad := strings.Repeat(" ", padding)
	return caption + "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clear() {
    fmt.Print("\033[2J") // clear the screen - works with vt100 terminals
    fmt.Print("\033[H") // move cursor to home
}

func main() {
    clear()
	if progressbar {
		prog := progress.New(progress.WithScaledGradient(fromcolour, tocolour))
		mod := model{
			progress:   prog,
			total_seconds: total_seconds,
			curseconds: curseconds,
		}
        if mode == Countup {
            mod.percent = 0
        } else if mode == Countdown {
            mod.percent = 100
        } else {
            panic("Unknown mode")
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
