package SecretGopher

// Game is the interface to the event handler
type Game struct {
	data gameData
	in   chan input
	out  chan Output
}

// GameState is a standalone type.
// GameState represents an instant of a game. All data contained in the struct is thread safe
// Depending on the Output type this struct is in, some values may be missing
type GameState struct {
	ElectionTracker int8   // ElectionTracker cycles from 0 to 3
	FascistTracker  int8   // FascistTracker starts at 0 ( no cards ), ends at 6 ( 6 slots )
	LiberalTracker  int8   // LiberalTracker starts at 0 ( no cards ), ends at 5 ( 5 slots )
	President       int8   // President is the current President (elected or candidate)
	Chancellor      int8   // Chancellor is the current Chancellor (elected or candidate)
	Roles           []Role // Roles is an array that maps a player's index to his role
	Votes           []Vote // Votes saves the votes for each player this round
	Killed          []int8 // Killed is a set that memorizes the ids of dead players
	Limited         []int8 // Limited is a set that memorizes the ids of limited players
}

// NewGame creates a game structure and subscribes a goroutine to listen to the events for the game
func NewGame() Game {
	G := Game{
		data: gameData{
			state:   waitingPlayers,
			players: 0,
			//deck:          // initialized later
			president:  NotSet,
			chancellor: NotSet,
			//roles:         // initialized later
			nextPresident: NotSet,
			oldGov:        make([]int8, 2),
			killed:        make([]int8, 2),
			//investigated:  // initialized later
			//votes:         // initialized later
			voted: 0,
			//policyChoice:  // initialized later
			eTracker: 0,
			fTracker: 0,
			lTracker: 0,
		},
	}
	G.subscribeHandler()
	return G
}

func (g *Game) Start() Output {
	g.in <- input{
		gameData: &g.data,
		event:    start{},
	}
	return <-g.out
}

func (g *Game) AddPlayer() Output {
	g.in <- input{
		gameData: &g.data,
		event:    addPlayer{},
	}
	return <-g.out
}

func (g *Game) Vote(c int8, v Vote) Output {
	g.in <- input{
		gameData: &g.data,
		event:    playerVote{Caller: c, Vote: v},
	}
	return <-g.out
}

func (g *Game) MakeChancellor(c, p int8) Output {
	g.in <- input{
		gameData: &g.data,
		event:    makeChancellor{Caller: c, Proposal: p},
	}
	return <-g.out
}

func (g *Game) PolicyDiscard(c int8, s uint8) Output {
	g.in <- input{
		gameData: &g.data,
		event:    policyDiscard{Caller: c, Selection: s},
	}
	return <-g.out
}

func (g *Game) SpecialPower(c int8, p SpecialPowers, s int8) Output {
	g.in <- input{
		gameData: &g.data,
		event:    specialPower{Caller: c, Power: p, Selection: s},
	}
	return <-g.out
}
