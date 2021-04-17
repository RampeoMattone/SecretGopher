package SecretGopher

// Game is the interface to the event handler
type Game struct {
	in  chan<- Event
	out <-chan Output
}

// GameState is a standalone type.
// GameState represents an instant of a game. All data contained in the struct is thread safe
// Depending on the Output type this struct is in, some values may be missing
type GameState struct {
	ElectionTracker int8   // ElectionTracker cycles from 0 to 3
	FascistTracker  int8   // FascistTracker starts at 0 ( no cards ), ends at 6 ( 6 slots )
	LiberalTracker  int8   // LiberalTracker starts at 0 ( no cards ), ends at 5 ( 5 slots )
	President       int8   // President is the current President (elected or candidate)
	Chancellor      int8   // Chancellor is the current Chancellor (elected or candidate)
	Roles           []Role // Roles is an array that maps a player's index to his role
	Votes           []Vote // Votes saves the votes for each player this round
	Killed          []int8 // Killed is a set that memorizes the ids of dead players
	Limited         []int8 // Limited is a set that memorizes the ids of limited players
}
