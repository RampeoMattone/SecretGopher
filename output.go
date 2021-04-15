package SecretGopher

type (
	// Output is a void interface. it's only used to simplify reading code
	Output interface{}

	// Error contains error information
	Error struct{ Err interface{} }

	// WrongPhase is an Error type.
	// WrongPhase means the event was sent at the wrong time
	WrongPhase struct{}

	// GameFull is an Error type.
	// GameFull means the number of players cannot grow anymore
	GameFull struct{}

	// Unauthorized is an Error type.
	// Unauthorized means the event was sent by the wrong authority (i.e. the wrong player)
	Unauthorized struct{}

	// Invalid is an Error type.
	// Invalid means the event was sent and contained Invalid data
	Invalid struct{}

	// Ok reports a successful interaction and contains information about how the game reacted
	Ok struct{ Info interface{} }

	// General is an Ok type.
	// General means the state of the game did not change
	General struct{}

	// PlayerRegistered is an Ok type.
	// PlayerRegistered means a player was registered.
	// The value associated with this type is the player's id
	PlayerRegistered int8

	// GameStart is an Ok type.
	// GameStart means the game has started.
	// the type contains the first president to be up for election and the roles that have been given to each player
	GameStart struct {
		President int8   // current President (elected or candidate)
		Roles     []Role // array that maps a player's index to his role
	}

	// NextPresident is an Ok type.
	// NextPresident that a round ended and a new president candidate was selected.
	// The new president's id is in the 'President' field.
	// NextPresident also carries a pointer to a GameState
	NextPresident struct {
		President int8
		State     *GameState
	}

	// ElectionResult is an Ok type.
	// ElectionResult means the voting phase has ended and its result is carried in the 'Result' field
	// ElectionResult also carries a pointer to a GameState
	ElectionResult struct {
		Result bool
		State  *GameState
	}

	// PolicyDiscardOk is an Ok type.
	// PolicyDiscardOk means the PolicyDiscard was successful.
	// PolicyDiscardOk also carries a pointer to a GameState
	PolicyDiscardOk *GameState

	// GameState is a standalone type.
	// GameState represents an instant of a game. All data contained in the struct is thread safe
	// Depending on the Output type this struct is in, some values may be missing
	GameState struct {
		ElectionTracker int8     // cycles from 0 to 3
		FascistTracker  int8     // starts at 0 ( no cards ), ends at 6 ( 6 slots )
		LiberalTracker  int8     // starts at 0 ( no cards ), ends at 5 ( 5 slots )
		President       int8     // current President (elected or candidate)
		Chancellor      int8     // current President (elected or candidate)
		Hand            []Policy // policy cards that the president or chancellor have to filter and choose from
		Roles           []Role   // array that maps a player's index to his role
	}
)
