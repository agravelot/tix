package main

import (
	"fmt"

	"github.com/agravelot/tix/project"
	"github.com/agravelot/tix/source/kitty"
	"github.com/pkg/errors"

	"github.com/rivo/tview"
)

type ProjectSource interface {
	ListProjects() ([]project.Project, error)
	StartProject(project project.Project) error
}

// Try to use bubbletea
// https://github.com/charmbracelet/bubbletea/tree/master
func main() {
	// r := tmux.TmuxSource{}
	r := kitty.KittySource{
		// TODO make it configurable
		ConfigPath:     "./config",
		RemotePassword: "my passphrase",
	}
	projects, err := r.ListProjects()
	if err != nil {
		panic(fmt.Errorf("unable listing projects: %w", err))
	}

	app := tview.NewApplication()

	list := tview.NewList()
	list.SetBorder(true).SetTitle(" tix ")

	// TODO FindItems
	for _, project := range projects {
		list.AddItem(project.Line(), "", 'i', nil)
	}

	list.SetSelectedFunc(func(index int, s string, sq string, _ rune) {
		p := &projects[index]
		p.Selected = !p.Selected
		list.SetItemText(index, p.Line(), "")

		for _, p := range projects {

			if !p.Selected {
				continue
			}

			err := r.StartProject(p)
			if err != nil {
				panic(errors.Wrap(err, "unable starting project: "))
			}
		}
	})

	list.SetDoneFunc(func() {
		for _, p := range projects {

			if !p.Selected {
				continue
			}

			err := r.StartProject(p)
			if err != nil {
				panic(errors.Wrap(err, "unable starting project: "))
			}
		}
		app.Stop()
	})

	button := tview.NewButton("Hit Enter to close").SetSelectedFunc(func() {
		app.Stop()
	})

	flexButtons := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(button, 0, 1, false).AddItem(button, 0, 1, false)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 1, true).
		AddItem(flexButtons, 2, 1, false)

	err = app.SetRoot(flex, true).SetFocus(flex).Run()
	if err != nil {
		panic(err)
	}
}
