package SecretGopher

type (
	input struct {
		*gameData
		event
	}

	event interface{} // event is a void interface. it's only used to simplify reading code

	// addPlayer is an event type.
	// addPlayer requests that the number of players be increased by one
	addPlayer struct{}

	// start is an event type.
	// start requests that the game starts
	start struct{}

	// makeChancellor is an event type.
	// makeChancellor requests that player 'Proposal' is made chancellor
	// under Caller's presidency
	makeChancellor struct {
		Caller   int8
		Proposal int8
	}

	// playerVote is an event type.
	// playerVote says that player 'Caller' has voted 'Vote' on either an election or a veto
	playerVote struct {
		Caller int8
		Vote   Vote
	}

	// policyDiscard is an event type.
	// policyDiscard says that player 'Caller' has decided to discard the policy card identified
	// by Selection
	policyDiscard struct {
		Caller    int8
		Selection uint8
	}

	// specialPower is an event type.
	// specialPower says the Power being used by 'Caller' under
	// the Power 'field' and the eventual selection of entity in the 'Selection' field
	specialPower struct {
		Caller    int8
		Power     SpecialPowers
		Selection int8
	}
)
