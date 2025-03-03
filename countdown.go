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
	"path"
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
	Reset    = "\033[0m"
	Red      = "\033[31m"
	Green    = "\033[32m"
	Yellow   = "\033[33m"
	Blue     = "\033[34m"
	Purple   = "\033[35m"
	Cyan     = "\033[36m"
	Gray     = "\033[37m"
	White    = "\033[97m"
	width    = 70
)

var (
	count_duration     time.Duration
	remaining_duration time.Duration
	seconds            int64
	minutes            int64
	hours              int64
	description        string
	until              string
	progressbar        bool
	fromcolour         = "#0000FF"
	tocolour           = "#FF0000"
	helpStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	endtime            time.Time
	now                time.Time
	debug              bool
	logfile            *os.File = nil
	logpath            string
)

func init() {
	flag.Int64Var(&seconds, "seconds", 0, "Number of seconds to count down")
	flag.Int64Var(&minutes, "minutes", 0, "Number of minutes to count down")
	flag.Int64Var(&hours, "hours", 0, "Number of hours to count down")
	flag.BoolVar(&progressbar, "prog", true, "Display an in-colour progress bar")
	flag.StringVar(&fromcolour, "fromcolour", "#0000FF", "Left-hand colour of gradient")
	flag.StringVar(&tocolour, "tocolour", "#FF0000", "Right-hand colour of gradient")
	flag.StringVar(&description, "description", "", "Description of what happens when the count is done")
	flag.StringVar(&until, "until", "", "Countup/down until time (HH:MM:SS)")
	flag.BoolVar(&debug, "debug", false, "Enable debug logging to $HOME/countdown.log")
	flag.Parse()

	var err error
	if debug {
		logpath = path.Join(os.Getenv("HOME"), "countdown.log")
		logfile, err = os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		debuglog("Starting count")
	}

	count_duration = time.Duration(int64(time.Second) * (seconds + minutes*60 + hours*3600))
	remaining_duration = count_duration

	if until != "" {
		fmt.Fprintf(os.Stderr, "I'm sorry, but --until functionality is still in development.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if remaining_duration == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// compute endtime so we always have a reference of when we are done
	now = time.Now()
	endtime = now.Add(count_duration)

	debuglog("now is %s", now)
	debuglog("endtime is %s", endtime)
}

func debuglog(format string, args ...interface{}) {
	if logfile == nil {
		return
	} else {
		fmt.Fprintf(logfile, format, args...)
		fmt.Fprintf(logfile, "\n")
	}
}

func telltime() {
	fmt.Printf("\r")
	now := time.Now().UTC()
	year := now.Local().Year()
	next_year := year + 1
	newyears := time.Date(next_year, time.January, 1, 0, 0, 0, 0, time.Local)
	diff := newyears.Sub(now)
	printstring := fmt.Sprintf("Countdown: %s until %d", mlib.Duration2Human(diff, false), next_year)
	finalstring := printstring[:]
	if len(printstring) > width {
		finalstring = printstring[:width]
	}
	format := fmt.Sprintf("%%s%%%ds%%s", width)
	fmt.Printf(format, Purple, finalstring, Reset)
}

type tickMsg time.Time

type model struct {
	percent            float64
	remaining_duration time.Duration
	progress           progress.Model
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
		// tickMsg should be the current time, no?
		now = time.Time(msg)

		m.remaining_duration = endtime.Sub(now)
		remaining_duration = m.remaining_duration

		debuglog("now is %v, endtime is %v", now, endtime)

		if now.Before(endtime) {
			debuglog("now is < endtime, we are still counting")
			// Then we're still counting.
			m.percent = float64(m.remaining_duration) / float64(count_duration)
			debuglog("remaining %d, count %d", m.remaining_duration, count_duration)
			debuglog("percent is %f", m.percent)
		} else {
			// We are done.
			debuglog("we are done")
			m.percent = 1.0
			return m, tea.Quit
		}

		return m, tickCmd()

	default:
		return m, nil
	}
}

func (m model) View() string {
	now = time.Now()
	current_duration := endtime.Sub(now)
	human_current := mlib.Duration2Human(current_duration, true)
	//human_total := mlib.Duration2Human(m.count_duration, true)
	caption := ""

	extra := ""
	if description != "" {
		extra = " until " + description
	}
	caption = fmt.Sprintf("Time remaining%s: %s", extra, human_current)
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
	fmt.Print("\033[H")  // move cursor to home
}

func main() {
	clear()
	if progressbar {
		prog := progress.New(progress.WithScaledGradient(fromcolour, tocolour))
		mod := model{
			progress:           prog,
			remaining_duration: remaining_duration,
			percent:            0.0,
		}
		mod.percent = 0

		if _, err := tea.NewProgram(mod).Run(); err != nil {
			fmt.Println("Oh no!", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Here I will put a plain old ASCII countdown")
		os.Exit(1)
	}
}
