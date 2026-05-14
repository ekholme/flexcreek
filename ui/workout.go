package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ekholme/flexcreek"
)

const (
	stateWorkoutList sessionState = iota
	stateCreateWorkout
	stateViewWorkout
)

// defining interaces that the workout model requires
type WorkoutProvider interface {
	GetLatestWorkouts(ctx context.Context, n int, userID int) ([]*flexcreek.Workout, error)
	GetWorkoutByID(ctx context.Context, id int, userID int) (*flexcreek.Workout, error)
}

type WorkoutCreator interface {
	CreateWorkout(ctx context.Context, w *flexcreek.Workout) (int, error)
}

type WorkoutDeleter interface {
	DeleteWorkout(ctx context.Context, id int) error
}

type WorkoutStore interface {
	WorkoutProvider
	WorkoutCreator
	WorkoutDeleter
}

// handles all interactions with the workout model
type WorkoutModel struct {
	store           WorkoutStore
	list            list.Model
	input           textinput.Model
	state           sessionState
	loading         bool
	err             error
	selectedUserID  int //i think this is the right way to handle this for now?
	listLength      int
	selectedWorkout *flexcreek.Workout
}

func NewWorkoutModel(s WorkoutStore, userID int, listLength int) WorkoutModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a Workout"

	//add an entry in the help keybinds to create a new workout
	var createWorkoutKey = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new workout"),
	)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			createWorkoutKey,
		}
	}

	//text input stuff
	ti := textinput.New()
	ti.Placeholder = "New Workout Short Description..."
	ti.Focus()

	return WorkoutModel{
		store:          s,
		list:           l,
		input:          ti,
		state:          stateWorkoutList,
		loading:        true,
		selectedUserID: userID,
		listLength:     listLength,
	}
}

// a command to fetch the latest workouts for a given user from the database
// again, this is wrapped in a command so it's non-blocking
func fetchLatestWorkoutsCmd(s WorkoutStore, n int, userID int) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		workouts, err := s.GetLatestWorkouts(ctx, n, userID)
		if err != nil {
			return err
		}

		return workoutsLoadedMsg{workouts}
	}
}

func createWorkoutCmd(s WorkoutStore, w *flexcreek.Workout) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		_, err := s.CreateWorkout(ctx, w)
		if err != nil {
			return err
		}

		return workoutCreatedMsg{}
	}
}

// struct wrappers for messages
type workoutsLoadedMsg struct {
	workouts []*flexcreek.Workout
}

type workoutSelectedMsg struct {
	workout *flexcreek.Workout
}

type workoutCreatedMsg struct {
}

type workoutItem struct {
	flexcreek.Workout
}

func (i workoutItem) Title() string       { return i.ShortDescription }
func (i workoutItem) Description() string { return i.WorkoutDate.Format("2006-01-02") }
func (i workoutItem) FilterValue() string { return i.LongDescription }

// bubbletea model requirements
func (m WorkoutModel) Init() tea.Cmd {
	return fetchLatestWorkoutsCmd(m.store, m.listLength, m.selectedUserID)
}

func (m WorkoutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		m.loading = false
		m.err = msg
		return m, nil

	case workoutsLoadedMsg:
		m.loading = false
		items := make([]list.Item, len(msg.workouts))
		for i, w := range msg.workouts {
			items[i] = workoutItem{*w}
		}
		m.list.SetItems(items)

	case workoutCreatedMsg:
		m.input.Blur()
		m.input.Reset()
		return m, fetchLatestWorkoutsCmd(m.store, m.listLength, m.selectedUserID)

	case tea.WindowSizeMsg:
		switch m.state {
		case stateWorkoutList:
			return m.updateWorkoutList(msg)
		case stateCreateWorkout:
			return m.updateWorkoutForm(msg)
		}

	default:
		switch m.state {
		case stateWorkoutList:
			return m.updateWorkoutList(msg)
		case stateCreateWorkout:
			return m.updateWorkoutForm(msg)
		}

	}

	return m, cmd
}

func (m WorkoutModel) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error()
	}

	switch m.state {
	case stateCreateWorkout:
		return "\n Placeholder \n\n" +
			m.input.View() +
			"\n\n (esc to go back)"

	default:
		if m.loading {
			return " Loading workouts..."
		}

		return "\n" + m.list.View()
	}
}

// update helpers
func (m WorkoutModel) updateWorkoutList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.list.SetSize(size.Width, size.Height)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		//prevents workout selection in a filtering state
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "n":
			m.state = stateCreateWorkout
			m.input.Focus()
			return m, nil

		case "enter":
			if i, ok := m.list.SelectedItem().(workoutItem); ok {
				m.selectedWorkout = &i.Workout
				return m, func() tea.Msg { return workoutSelectedMsg{&i.Workout} }
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m WorkoutModel) updateWorkoutForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	//todo
	return nil, nil
}
