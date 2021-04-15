package SecretGopher

type (
	// Output is a void interface. it's only used to simplify reading code
	Output interface{}

	// Error contains error information
	Error struct{ e interface{} }

	// GameFull is an Error type. GameFull means the number of players cannot grow anymore
	GameFull struct{}

	// WrongPhase is an Error type. WrongPhase means the event was sent at the wrong time
	WrongPhase struct{}
)
