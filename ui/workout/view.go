package workout

import "fmt"

func (m Model) View() string {
	return fmt.Sprintf("Welcome, %s!", m.User.Username)
}
