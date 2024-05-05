package project

type Project struct {
	Name     string
	Selected bool
	Opened   bool
}

func (p Project) Line() string {
	color := ""
	selected := "[   ]"

	if p.Opened {
		color = "[green]"
	}

	if p.Selected {
		selected = "[ x ]"
	}

	return color + selected + " - " + p.Name
}
