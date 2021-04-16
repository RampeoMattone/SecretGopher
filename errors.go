package SecretGopher

type (
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
)
