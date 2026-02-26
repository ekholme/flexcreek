package ui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
	"github.com/ekholme/flexcreek/ui/userselect"
	"github.com/ekholme/flexcreek/ui/workout"
)

type sessionState int

const (
	userSelectView sessionState = iota
	workoutView
)

type MainModel struct {
	state sessionState

	userselect userselect.Model
	workout    workout.Model

	//services
	userService    flexcreek.UserService
	workoutService flexcreek.WorkoutService
}


func (m MainModel) Init() tea.Cmd {
	return m.userselect.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case userSelectView:
		if m.userselect.SelectedUser != nil {
			m.state = workoutView
			m.workout = workout.New(m.userselect.SelectedUser)
			return m, m.workout.Init()
		}
		newModel, newCmd := m.userselect.Update(msg)
		m.userselect = newModel.(userselect.Model)
		cmd = newCmd

	case workoutView:
		newModel, newCmd := m.workout.Update(msg)
		m.workout = newModel.(workout.Model)
		cmd = newCmd
	}
	return m, cmd
}

func (m MainModel) View() string {
	switch m.state {
	case workoutView:
		return m.workout.View()
	default:
		return m.userselect.View()
	}
}
