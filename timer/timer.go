package main

import (
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/bubbles/timer"
    "github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

var bigTextStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("5")).
    Background(lipgloss.Color("#FAFAFA"))

type model struct {
    timer           timer.Model
    originalTimeout time.Duration
    keymap          keymap
    help            help.Model
    quitting        bool
}

type keymap struct {
    start key.Binding
    stop  key.Binding
    reset key.Binding
    quit  key.Binding
}

func (m model) Init() tea.Cmd {
    return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case timer.TickMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd

    case timer.StartStopMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        m.keymap.stop.SetEnabled(m.timer.Running())
        m.keymap.start.SetEnabled(!m.timer.Running())
        return m, cmd

    case timer.TimeoutMsg:
        m.quitting = true
        return m, tea.Quit

    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keymap.quit):
            m.quitting = true
            return m, tea.Quit
        case key.Matches(msg, m.keymap.reset):
            m.timer = timer.NewWithInterval(m.originalTimeout, time.Millisecond*100)
            return m, m.timer.Init()
        case key.Matches(msg, m.keymap.start, m.keymap.stop):
            return m, m.timer.Toggle()
        }
    }

    return m, nil
}

func (m model) helpView() string {
    return "\n" + m.help.ShortHelpView([]key.Binding{
        m.keymap.start,
        m.keymap.stop,
        m.keymap.reset,
        m.keymap.quit,
    })
}

func (m model) View() string {
    s := bigTextStyle.Render(m.timer.View())
    if m.timer.Timedout() {
        s = bigTextStyle.Render("Time's up!")
    } else {
        s = "Timer initialized: "+ time.Now().Format("2006-01-02 15:04:05") + "\n" + "Time Remaining: " + s
    }
    s += "\n" + m.helpView()
    return s
}

func main() {
    var duration int
    flag.IntVar(&duration, "duration", 1, "Set the duration for the timer in minutes")
    flag.Parse()

    // Calculate timeout from duration
    timeout := time.Minute * time.Duration(duration)

    m := model{
        timer:           timer.NewWithInterval(timeout, time.Millisecond*100),
        originalTimeout: timeout,
        keymap: keymap{
            start: key.NewBinding(
                key.WithKeys("s"),
                key.WithHelp("s", "start/stop the timer"),
            ),
            stop: key.NewBinding(
                key.WithKeys("s"),
                key.WithHelp("s", "start/stop the timer"),
            ),
            reset: key.NewBinding(
                key.WithKeys("r"),
                key.WithHelp("r", "reset the timer"),
            ),
            quit: key.NewBinding(
                key.WithKeys("q", "ctrl+c"),
                key.WithHelp("q", "quit the program"),
            ),
        },
        help: help.New(),
    }

    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Println("Uh oh, we encountered an error:", err)
        os.Exit(1)
    }
}
