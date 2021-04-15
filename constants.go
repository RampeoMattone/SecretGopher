package SecretGopher

const (
	notSet int8 = -1
)

type state uint8 // state enumerates the possible round states

const (
	waitingPlayers      state = iota // waitingPlayers means the game is waiting a Start or an AddPlayer Event
	chancellorCandidacy              // chancellorCandidacy means the game is waiting a MakeChancellor Event
)
