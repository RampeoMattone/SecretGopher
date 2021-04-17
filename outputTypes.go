package SecretGopher

/*
Output events convention:
The library shall use types as a way to communicate a meaningful change in state.
The library shall also ship a copy of its state within those events, such that a state reconstruction is possible from the outside.
The library may omit a gamestate event and use a slimmer Ok type to signal that the state did not change meaningfully
*/

type (
	// Output is a void interface. it's only used to simplify reading code
	Output interface{}

	// Ok reports a successful interaction and contains information about how the game reacted
	Ok struct{ Info interface{} }

	// VoteRegistered is an Ok type.
	// VoteRegistered means the vote was valid but the state of the game did not change
	VoteRegistered struct{}

	// PlayerRegistered is an Ok type.
	// PlayerRegistered means a player was registered.
	// The value associated with this type is the player's id
	PlayerRegistered int8

	// GameStart is an Ok type.
	// GameStart means the game has started.
	// GameStart also carries a pointer to a GameState
	GameStart GameState

	// NextPresident is an Ok type.
	// NextPresident that a round ended and a new president candidate was selected.
	// NextPresident also carries a pointer to a GameState
	NextPresident GameState

	// ElectionStart is an Ok type.
	// ElectionStart that a round ended and a new president candidate was selected.
	// ElectionStart also carries a pointer to a GameState
	ElectionStart GameState

	// LegislationPresident is an Ok type.
	// LegislationPresident means the voting phase has ended successfully and the legislative session has started.
	// LegislationPresident also carries a pointer to a GameState
	LegislationPresident struct {
		Hand  []Policy
		State GameState
	}

	// LegislationChancellor is an Ok type.
	// LegislationChancellor means the chancellor has to select a policy to enact
	// LegislationChancellor also carries a pointer to a GameState
	LegislationChancellor struct {
		Hand  []Policy
		State GameState
	}

	// PolicyEnaction is an Ok type.
	// PolicyEnaction means nobody won, but the Policy got enacted.
	// If a special power has been activated, the SpecialPower field will let you know
	// PolicyEnaction also carries a pointer to a GameState
	PolicyEnaction struct {
		Enacted      Policy
		SpecialPower SpecialPowers
		State        GameState
	}

	// SpecialPowerFeedback is an Ok type.
	// SpecialPowerFeedback contains information about the power that has just been used.
	// The Feedback field will carry a Policy slice in response to a Peek power, or a Role value in response to an Investigate power
	SpecialPowerFeedback struct {
		Feedback interface{}
		State    GameState
	}

	// VetoRequest is an Ok type.
	// VetoRequest means a veto is possible and the handler is now waiting for one or more VetoResponse inputs.
	VetoRequest GameState

	// GameEnd is an Ok type.
	// GameEnd means a condition to end the game has been met. the reason for the ending is in the field 'Why'
	// GameEnd also carries a pointer to a GameState
	GameEnd struct {
		Why   GameEnding
		State GameState
	}
)
