package SecretGopher

// NewGame creates a game structure and subscribes a goroutine to listen to the events for the game
func NewGame() Game {
	in := make(chan Event)
	out := make(chan Output)
	go handleGame(in, out)
	return Game{
		In:  in,
		Out: out,
	}
}

// handleGame handles the game events
func handleGame(in <-chan Event, out chan<- Output) {
	// information about the game is kept in the handler to ensure thread safety
	var (
		state      state = waitingPlayers
		players    int8  = 0
		president  int8  = notSet
		chancellor int8  = notSet
		eTracker   int8  = 0
		fTracker   int8  = 0
		lTracker   int8  = 0
	)
	for {
		event := <-in
		switch event.(type) {
		case AddPlayer:
			// if the game is accepting players
			if state == waitingPlayers {
				if players < 10 {
					players++ // adds a player to the game
				} else {
					out <- Error{e: GameFull{}} // send out error
				}
			} else {
				out <- Error{e: WrongPhase{}} // send out error
			}
		case Start:
			// if the game was accepting players
			if state == waitingPlayers {
				// todo select a random player to be president
				// state is advanced
				state = chancellorCandidacy
			} else {
				out <- Error{e: WrongPhase{}} // send out error
			}
		}
	}
}
