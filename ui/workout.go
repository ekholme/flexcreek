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
	store          WorkoutStore
	list           list.Model
	input          textinput.Model
	state          sessionState
	loading        bool
	err            error
	selectedUserID int //i think this is the right way to handle this for now?
	listLength     int
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

type workoutCreatedMsg struct {
}

type workoutItem struct {
	flexcreek.Workout
}

//RESUME HERE
