package ui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/agravelot/tix/core"
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

type item struct {
	Ws core.Workspace
}

func (i item) FilterValue() string { return i.Ws.Name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		log.Println("could not assert item to item")
		return
	}

	appletsStatus := "🔴"
	if i.Ws.IsRunning() {
		appletsStatus = "🟢"
	}

	appletsStatus += " \ue7b0 \uf308"

	for _, a := range i.Ws.Applets {
		if a.IsRunning {
			appletsStatus += "🟢"
			continue
		}
		appletsStatus += "🔴"
	}

	str := fmt.Sprintf("%d. %s %s", index+1, i.Ws.Name, appletsStatus)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type uiApp struct {
	application         core.Application
	list                list.Model
	settingUpWorkspaces bool
	quitting            bool
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
			log.Println("settingUpWorkspaces : ", m.selected)
			m.settingUpWorkspaces = true

			for i, w := range m.application.Workspaces {
				_, isSelected := m.selected[i]
				// TODO Refactor without nested if
				if w.IsRunning() != isSelected {
					if isSelected {
						log.Printf("Setting up workspace %s", w.Name)
						log.Printf("Running commands %v", w.SetupCommands)
						err := w.Setup(m.application.Docker)
						if err != nil {
							log.Fatalln("error setting up workspace : ", err)
						}
					} else {
						log.Printf("Tearing down workspace %s", w.Name)
						log.Printf("Running commands %v", w.TeardownCommands)
						err := w.Teardown(m.application.Docker)
						if err != nil {
							log.Fatalln("error tearing down workspace : ", err)
						}
					}
				}
			}

			return m, tea.Quit

		case " ":
			m.toggle()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *uiApp) toggle() {
	s := m.list.SelectedItem()
	log.Println("Selected item : ", s)
	_, ok := m.selected[m.list.Index()]
	if ok {
		delete(m.selected, m.list.Index())
		return
	}
	m.selected[m.list.Index()] = struct{}{}
	log.Println("Selected : ", m.list.Index())
}

func (m uiApp) View() string {
	// TODO Correctly handle slice
	if m.settingUpWorkspaces && len(m.selected) != 0 {
		out := ""

		for k := range m.selected {
			log.Println("Selected : ", k)
			out += fmt.Sprintf("%v\n", m.application.Workspaces[k].Name)
		}

		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", out))
	}

	if m.quitting {
		return quitTextStyle.Render("Goodbye! 👋")
	}

	return "\n" + m.list.View()
}

// New starts the UI
func New(app core.Application) error {
	const defaultWidth = 20

	items := []list.Item{}
	s := make(map[int]struct{}, len(app.Workspaces))

	for i := range app.Workspaces {
		items = append(items, item{Ws: app.Workspaces[i]})
		if app.Workspaces[i].IsRunning() {
			s[i] = struct{}{}
		}
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "tix"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := uiApp{application: app, list: l, selected: s}

	_, err := tea.NewProgram(m).Run()
	return err
}
