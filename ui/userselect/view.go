package userselect

import "fmt"

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %s", m.err)
	}

	switch m.state {
	case creating:
		return m.input.View()
	default:
		return m.list.View()
	}
}