package SecretGopher

// Game is the interface to the event handler
type Game struct {
	In  chan<- Event
	Out <-chan Output
}

// GameState is a struct that holds minimal information about the game
type GameState struct {
	// Public Game Sate
	ElectionTracker int8 // cycles from 0 to 3
	FascistTracker  int8 // starts at 0 ( no cards ), ends at 6 ( 6 slots )
	LiberalTracker  int8 // starts at 0 ( no cards ), ends at 5 ( 5 slots )
	President       int8 // current President (elected or candidate)
	Chancellor      int8 // current President (elected or candidate)
}
