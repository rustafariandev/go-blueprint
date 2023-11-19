package inputoptions

import "fmt"

type Item struct {
	Name, Desc, Filter, Value string
	Selected                  bool
}

func (i *Item) ToggleSelected() {
	i.Selected = !i.Selected
}

func (i *Item) Title() string {
	if i.Selected {
		return fmt.Sprintf("[x] %s", i.Name)
	}

	return fmt.Sprintf("[ ] %s", i.Name)
}

func (i *Item) Description() string { return "    " + i.Desc }
func (i *Item) FilterValue() string {
	if i.Filter != "" {
		return i.Filter
	}

	return i.Name
}
