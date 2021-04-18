package SecretGopher

const (
	NotSet int8 = -1
)

type state uint8 // state enumerates the possible round states

const (
	waitingPlayers        state = iota // waitingPlayers means the game is waiting a Start or an AddPlayer Event
	chancellorCandidacy                // chancellorCandidacy means the game is waiting a MakeChancellor Event
	governmentElection                 // governmentElection means the game is waiting a PlayerVote Event
	presidentLegislation               // presidentLegislation means the game is waiting a PolicyDiscard Event from the president
	chancellorLegislation              // chancellorLegislation means the game is waiting a PolicyDiscard Event from the chancellor
	specialPeek                        // presidentLegislation means the game is waiting a PolicyDiscard Event from the president
	specialInvestigate                 // presidentLegislation means the game is waiting a PolicyDiscard Event from the president
	specialElection                    // presidentLegislation means the game is waiting a PolicyDiscard Event from the president
	specialExecution                   // presidentLegislation means the game is waiting a PolicyDiscard Event from the president
	vetoChancellor
	vetoPresident
	gameEnd
)

// Role is used to represent the role of a player
type Role int8

const (
	LiberalParty Role = iota
	FascistParty
	Hitler
)

// Vote is used to represent what a player may decide using the Ja or Nein cards that the board game uses
type Vote int8

const (
	NoVote Vote = iota // NoVote means the player still hasn't voted
	Ja                 // Ja means in favor
	Nein               // Nein means against
)

// Policy is used to represent a Policy card
type Policy bool

const (
	LiberalPolicy Policy = true  // LiberalPolicy means the policy is in favor of the liberal party
	FascistPolicy Policy = false // FascistPolicy means the policy is in favor of the fascist party
)

type SpecialPowers int8

const (
	Nothing SpecialPowers = iota
	Peek
	Investigate
	Election
	Execution
)

// powersTable is used to figure out what power needs to be activated for a specific round
var powersTable = [3][6]SpecialPowers{
	{Nothing, Nothing, Peek, Execution, Execution, Nothing},
	{Nothing, Investigate, Election, Execution, Execution, Nothing},
	{Investigate, Investigate, Election, Execution, Execution, Nothing},
}

// GameEnding is used to signal if the game ended and how
type GameEnding int8

const (
	StillRunning        GameEnding = iota
	LiberalPolicyWin               // LiberalPolicyWin means 5 liberal policies have been enacted
	LiberalExecutionWin            // LiberalExecutionWin means hitler was killed
	FascistPolicyWin               // FascistPolicyWin means 6 fascist policies have been enacted
	FascistElectionWin             // FascistElectionWin means hitler was elected as chancellor
)