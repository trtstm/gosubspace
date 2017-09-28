package main

// GameState Is a state the game can be in. E.g Loading,Playing,...
type GameState interface {
	// Run Runs the current state. Return true if this state is done and NextState can be called.
	Run() bool
	// NextState returns the next state to run.
	NextState() GameState

	// Name gets the name of the state. Used mainly for debugging.
	Name() string
}

type LoadingState struct {
}

func (s *LoadingState) Run() bool {
	return false
}

func (s *LoadingState) NextState() GameState {
	return nil
}

func (s *LoadingState) Name() string {
	return "LoadingState"
}
