package SecretGopher

const (
	NotSet int8 = -1
)

type state uint8 // state enumerates the possible round states

const (
	waitingPlayers      state = iota // waitingPlayers means the game is waiting a Start or an AddPlayer Event
	chancellorCandidacy              // chancellorCandidacy means the game is waiting a MakeChancellor Event
	governmentElection               // governmentElection means the game is waiting a GovernmentVote Event
	presidentTurn                    // presidentTurn means the game is waiting a PolicyDiscard Event from the president
	chancellorTurn                   // chancellorTurn means the game is waiting a PolicyDiscard Event from the chancellor

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
