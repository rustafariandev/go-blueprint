package inputoptions

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Doc       lipgloss.Style
	Title     lipgloss.Style
	StatusMsg lipgloss.Style
}

var AppStyles = DefaultStyles()

func DefaultStyles() Styles {
	return Styles{
		Doc: lipgloss.NewStyle().Padding(1, 2),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Padding(0, 1),
		StatusMsg: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}),
	}
}

var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("21")).Bold(true)

const logo = `

 ____  _                       _       _   
|  _ \| |                     (_)     | |  
| |_) | |_   _  ___ _ __  _ __ _ _ __ | |_ 
|  _ <| | | | |/ _ \ '_ \| '__| | '_ \| __|
| |_) | | |_| |  __/ |_) | |  | | | | | |_ 
|____/|_|\__,_|\___| .__/|_|  |_|_| |_|\__|
                                  | |
                                  |_|

`

var (
	docStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

type Model struct {
	Size     tea.WindowSizeMsg
	input    textinput.Model
	list     list.Model
	header   string
	output   *Output
	quit     *bool
	showList bool
	skipList bool
}

func (m Model) SetQuit() {
	*m.quit = true
}

func (m Model) Quit() bool {
	return *m.quit
}

type ModelOptions struct {
	Items      []*Item
	Header     string
	ListHeader string
	ShowList   bool
	SkipList   bool
}

func ValidateName(string) error {
	return nil
}

func NewModel(options ModelOptions) Model {
	listItems := ToListItems(options.Items)
	focusedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	input := textinput.New()
	input.CharLimit = 256
	input.Cursor.Style = focusedStyle.Copy()
	input.PromptStyle = focusedStyle
	input.TextStyle = focusedStyle
	input.Focus()
	output := &Output{}
	m := Model{list: list.New(listItems, NewItemDelegate(&output.Framework), 0, 0), input: input}
	m.quit = new(bool)
	m.list.Title = options.ListHeader
	m.list.Styles.Title = titleStyle
	m.output = output
	m.header = options.Header
	m.skipList = options.SkipList
	m.showList = options.ShowList
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.SetQuit()
			return m, tea.Quit
		}
	}

	if m.showList {
		if m.skipList {
			return m, tea.Quit
		}

		return m.listUpdate(msg)
	}

	return m.inputUpdate(msg)
}

func (m Model) inputUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			s := m.input.Value()
			if s != "" || m.input.Err == nil {
				m.input.Blur()
				m.showList = true
				m.output.Name = m.input.Value()
				return m, nil
			}
		}
	}

	i, cmd := m.input.Update(msg)
	m.input = i
	return m, tea.Batch(cmd)
}

func (m Model) listUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyRunes:
			if msg.String() == "q" {
				m.SetQuit()
				return m, tea.Quit
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showList {
		if !m.skipList {
			return m.listView()
		}
	}

	return m.inputView()

}

func (m Model) inputView() string {
	var b strings.Builder
	b.WriteString(logoStyle.Render(logo) + "\n")
	b.WriteString(m.header)
	b.WriteString("\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	if m.input.Err != nil {
		b.WriteString(errorStyle.Render(m.input.Err.Error()))
	}

	return b.String()
}

func (m Model) listView() string {
	return docStyle.Render(m.list.View())
}

func NewItemDelegate(out *string) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	choose := key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "choose"),
	)

	enter := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
	)

	d.SetSpacing(0)

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		item, ok := m.SelectedItem().(*Item)
		if !ok {
			return nil
		}

		title := item.Name

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, choose):
				items := m.Items()
				index := m.Index()
				for i := 0; i < index; i++ {
					if item, ok := items[i].(*Item); ok {
						item.selected = false
					}
				}

				item.ToggleSelected()
				for i := index + 1; i < len(items); i++ {
					if item, ok := items[i].(*Item); ok {
						item.selected = false
					}
				}

				if item.Selected() {
					*out = item.Value
					return m.NewStatusMessage(statusMessageStyle("You choose " + title))
				} else {
					*out = ""
				}

			}

		}

		return nil
	}

	help := []key.Binding{choose, enter}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

func ToListItems[T list.Item](in []T) []list.Item {
	items := make([]list.Item, len(in))
	for i, item := range in {
		items[i] = item
	}

	return items
}

type Output struct {
	Name      string
	Framework string
}

func (m Model) GetOutput() *Output {
	return m.output
}
