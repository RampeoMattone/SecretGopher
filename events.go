package SecretGopher

type (
	Event interface{} // Event is a void interface. it's only used to simplify reading code

	// AddPlayer is an Event type.
	// AddPlayer requests that the number of players be increased by one
	AddPlayer struct{}

	// Start is an Event type.
	// Start requests that the game starts
	Start struct{}

	// MakeChancellor is an Event type.
	// MakeChancellor requests that player 'C' is made chancellor
	MakeChancellor struct{ C int8 }
)
