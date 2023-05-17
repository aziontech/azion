package insert

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"strings"
)

type (
	errMsg error
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)                                   // nolint:all
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")) // nolint:all
)

type item string // nolint:all

func (i item) FilterValue() string { return "" } // nolint:all

type itemDelegate struct{} // nolint:all

func (d itemDelegate) Height() int                               { return 1 }   // nolint:all
func (d itemDelegate) Spacing() int                              { return 0 }   // nolint:all
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil } // nolint:all
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) { // nolint:all
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	textInput textinput.Model
	err       error
} // nolint:all

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	NameProject = m.textInput.Value()
	return m, cmd
}

var NameProject string = "undefined name"

func (m model) View() string {
	return fmt.Sprintf("What name project?\n\n%s\n\n%s", m.textInput.View(), "\n")
}

func Insert() {
	ti := textinput.New()
	ti.Placeholder = "name_project"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		textInput: ti,
		err:       nil,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(err)
	}
}