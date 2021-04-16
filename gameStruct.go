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
	ElectionTracker int8     // cycles from 0 to 3
	FascistTracker  int8     // starts at 0 ( no cards ), ends at 6 ( 6 slots )
	LiberalTracker  int8     // starts at 0 ( no cards ), ends at 5 ( 5 slots )
	President       int8     // current President (elected or candidate)
	Chancellor      int8     // current President (elected or candidate)
	Roles           []Role   // array that maps a player's index to his role
	Killed			Set		 // Set that memorizes the ids of dead players
}

// todo killed is not shared normally, find a way to do it
// todo gamehandler must become a method and repeated code needs to be fixed
// todo add veto check