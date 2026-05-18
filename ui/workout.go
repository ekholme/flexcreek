package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
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

type WorkoutModelInputs struct {
	ShortDescriptionInput textinput.Model
	LongDescriptionInput  textarea.Model
	WorkoutDateInput      textinput.Model
}

// handles all interactions with the workout model
type WorkoutModel struct {
	store           WorkoutStore
	list            list.Model
	inputs          WorkoutModelInputs
	inputFocusIndex int
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

	//short description init
	sdi := textinput.New()
	sdi.Placeholder = "Short Description (e.g. KB ABC)"
	sdi.Focus()

	//long description init
	ldi := textarea.New()
	ldi.Placeholder = "Long Description (e.g. 20 min AMRAP...)"

	wdi := textinput.New()
	wdi.Placeholder = "Workout Date (YYYY-MM-DD)"
	wdi.CharLimit = 10

	wmi := WorkoutModelInputs{
		ShortDescriptionInput: sdi,
		LongDescriptionInput:  ldi,
		WorkoutDateInput:      wdi,
	}

	return WorkoutModel{
		store:          s,
		list:           l,
		inputs:         wmi,
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
		// Reset form and go back to list
		m.state = stateWorkoutList
		m.loading = true
		m.inputs.ShortDescriptionInput.Reset()
		m.inputs.LongDescriptionInput.Reset()
		m.inputs.WorkoutDateInput.Reset()
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
		case stateViewWorkout:
			return m.updateViewWorkout(msg)
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
		return m.viewWorkoutForm()

	case stateViewWorkout:
		if m.selectedWorkout == nil {
			return "Error: No workout selected."
		}
		return "\n" + m.selectedWorkout.ShortDescription + "\n\n" +
			"Date: " + m.selectedWorkout.WorkoutDate.Format("2006-01-02") + "\n\n" +
			m.selectedWorkout.LongDescription + "\n\n" +
			"(esc to go back)"
	default:
		if m.loading {
			return " Loading workouts..."
		}

		return "\n" + m.list.View()
	}
}

// view helper for the create workout form
func (m WorkoutModel) viewWorkoutForm() string {
	return "\n Create New Workout \n\n" +
		m.inputs.ShortDescriptionInput.View() + "\n\n" +
		m.inputs.LongDescriptionInput.View() + "\n\n" +
		m.inputs.WorkoutDateInput.View() + "\n\n" +
		"(esc to go back)"
}

func (m WorkoutModel) updateViewWorkout(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		m.state = stateWorkoutList
	}
	return m, nil
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
			m.inputs.ShortDescriptionInput.Focus()
			return m, nil

		case "enter":
			if i, ok := m.list.SelectedItem().(workoutItem); ok {
				m.state = stateViewWorkout
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = stateWorkoutList
			return m, nil

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button is focused?
			// If so, create the workout.
			if s == "enter" && m.inputFocusIndex == 2 { // 2 is the last input index
				dateStr := m.inputs.WorkoutDateInput.Value()
				t, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					// For now, we'll just use the current time if parsing fails.
					// A better approach would be to show a validation error to the user.
					t = time.Now()
				}

				w := flexcreek.Workout{
					UserID:           m.selectedUserID,
					ShortDescription: m.inputs.ShortDescriptionInput.Value(),
					LongDescription:  m.inputs.LongDescriptionInput.Value(),
					WorkoutDate:      t,
				}
				m.loading = true
				return m, createWorkoutCmd(m.store, &w)
			}

			// Cycle focus
			if s == "up" || s == "shift+tab" || (s == "enter" && m.inputFocusIndex == 1) { // Special case for textarea
				m.inputFocusIndex--
			} else {
				m.inputFocusIndex++
			}

			// Wrap focus
			if m.inputFocusIndex > 2 {
				m.inputFocusIndex = 0
			} else if m.inputFocusIndex < 0 {
				m.inputFocusIndex = 2
			}

			// Blur all inputs
			m.inputs.ShortDescriptionInput.Blur()
			m.inputs.LongDescriptionInput.Blur()
			m.inputs.WorkoutDateInput.Blur()

			// Focus the correct input
			switch m.inputFocusIndex {
			case 0:
				cmd = m.inputs.ShortDescriptionInput.Focus()
			case 1:
				cmd = m.inputs.LongDescriptionInput.Focus()
			case 2:
				cmd = m.inputs.WorkoutDateInput.Focus()
			}

			return m, cmd
		}
	}

	// Handle character input and blinking for the focused field
	cmd = m.updateFocusedInput(msg)

	return m, cmd
}

// helper to update the currently focused input field
func (m *WorkoutModel) updateFocusedInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.inputFocusIndex {
	case 0:
		m.inputs.ShortDescriptionInput, cmd = m.inputs.ShortDescriptionInput.Update(msg)
	case 1:
		m.inputs.LongDescriptionInput, cmd = m.inputs.LongDescriptionInput.Update(msg)
	case 2:
		m.inputs.WorkoutDateInput, cmd = m.inputs.WorkoutDateInput.Update(msg)
	}
	return cmd
}
