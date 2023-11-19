package inputoptions

import "fmt"

type Item struct {
	Name, Desc, Filter, Value string
	selected                  bool
}

func (i *Item) ToggleSelected() {
	i.selected = !i.selected
}

func (i *Item) Selected() bool {
	return i.selected
}

func (i *Item) Title() string {
	if i.selected {
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
