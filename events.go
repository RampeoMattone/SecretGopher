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
	// MakeChancellor requests that player 'Proposal' is made chancellor
	// under Caller's presidency
	MakeChancellor struct {
		Caller   int8
		Proposal int8
	}

	// GovernmentVote is an Event type.
	// GovernmentVote says that player 'Caller' has voted 'Vote'
	GovernmentVote struct {
		Caller int8
		Vote   Vote
	}

	// PolicyDiscard is an Event type.
	// PolicyDiscard says that player 'Caller' has decided to discard the policy card identified
	// by Selection
	PolicyDiscard struct {
		Caller    int8
		Selection uint8
	}
)
