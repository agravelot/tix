package ui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/agravelot/tix/app"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
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

type uiApp struct {
	cfg                 app.Config
	list                list.Model
	settingUpWorkspaces bool
	quitting            bool
	cursor              int
	selected            map[int]struct{}
}

func (m uiApp) Init() tea.Cmd {
	return nil
}

func (m uiApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			// i, ok := m.list.SelectedItem().(item)
			// if ok {
			// 	m.choice = string(i)
			// }
			log.Println("settingUpWorkspaces : ", m.selected)
			m.selected[m.cursor] = struct{}{}
			m.settingUpWorkspaces = true

			if len(m.selected) != 0 {
				for k := range m.selected {
					workspace := m.cfg.Workspaces[k]
					log.Printf("Setting up workspace %s", workspace.Name)
					log.Printf("Running commands %v", workspace.SetupCommands)
					workspace.Setup()
				}
			}

			return m, tea.Quit

		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m uiApp) View() string {
	// TODO Correctly handle slice
	if m.settingUpWorkspaces && len(m.selected) != 0 {
		out := ""

		for k := range m.selected {
			log.Println("Selected : ", k)
			out += fmt.Sprintf("%v\n", m.cfg.Workspaces[k].Name)
		}

		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", out))
	}

	if m.quitting {
		return quitTextStyle.Render("Goodbye! 👋")
	}

	return "\n" + m.list.View()
}

// New starts the UI
func New(cfg app.Config) error {
	const defaultWidth = 20

	items := []list.Item{}

	for _, w := range cfg.Workspaces {
		items = append(items, item(w.Name))
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "tix"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := uiApp{cfg: cfg, list: l, selected: map[int]struct{}{}}

	_, err := tea.NewProgram(m).Run()
	return err
}
