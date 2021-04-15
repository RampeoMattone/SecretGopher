package SecretGopher

import (
	"math/rand"
	"time"
)

// init makes sure the seed is different at every re-run of the library
// todo move this out of the library
func init() {
	rand.Seed(time.Now().Unix())
}

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

// handleGame handles the game events.
func handleGame(in <-chan Event, out chan<- Output) {
	// information about the game is kept in the handler to ensure thread safety
	var (
		state         state    = waitingPlayers
		players       int8     = 0
		deck          Deck     // initialized with Start
		president     int8     = NotSet
		chancellor    int8     = NotSet
		roles         []Role   // initialized with Start
		nextPresident int8     = NotSet
		oldGov        Set      = make(Set, 2)
		votes         []Vote   // initialized with Start
		voteCounter   int8     = 0
		policyChoice  []Policy // initialized when entering presidentTurn
		eTracker      int8     = 0
		fTracker      int8     = 0
		lTracker      int8     = 0
	)
	for {
		event := <-in
		switch event.(type) {
		case AddPlayer:
			// if the game is accepting players
			if state == waitingPlayers {
				if players < 10 {
					players++                                  // adds a player to the game
					out <- Ok{Info: PlayerRegistered(players)} // say the player was registered under the player number
				} else {
					out <- Error{Err: GameFull{}} // send out error
				}
			} else {
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case Start:
			// if the game was accepting players
			if state == waitingPlayers {
				if players >= 5 {
					roles = make([]Role, players) // initialize roles to the proper size
					votes = make([]Vote, players) // initialize votes to the proper size
					deck = newDeck()              // initialize deck and shuffle it

					roles[rand.Intn(int(players))] = Hitler // set one player as Hitler
					var nF int                              // number of fascists based on the lobby size
					switch players {
					case 5, 6:
						nF = 1
					case 7, 8:
						nF = 2
					case 9, 10:
						nF = 3
					}
					// assign nF FascistParty roles randomly
					for i := 0; i < nF; {
						// extract a player
						// if the role for that player is not FascistParty or Hitler, set him as FascistParty
						// and increase the counter
						if r := rand.Intn(int(players)); roles[r] == LiberalParty {
							roles[r] = FascistParty
							i++
						}
					}
					// the first player to be president is random
					president = int8(rand.Intn(int(players)))
					// set the next president in line
					nextPresident = (president + 1) % players

					state = chancellorCandidacy // after a president is selected, a chancellor needs to be selected

					out <- Ok{Info: GameStart{
						President: president,
						Roles:     append([]Role{}, roles...), // clone roles
					}} // tell the caller the game has started
				}
			} else {
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case MakeChancellor:
			// if the game was accepting players
			if state == chancellorCandidacy {
				e := event.(MakeChancellor)
				if e.Caller == president {
					if !oldGov.Has(e.Proposal) {
						chancellor = e.Proposal
						state = governmentElection
						out <- Ok{Info: General{}} // say the chancellor registration was successful
					} else {
						out <- Error{Err: Invalid{}} // send out error
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			} else {
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case GovernmentVote:
			// if the game is waiting for votes on the election
			if state == governmentElection {
				e := event.(GovernmentVote)
				// check that the vote is valid
				if v := e.Vote; v == Ja || v == Nein {
					// if the user hasn't voted yet
					if votes[e.Caller] == NoVote {
						votes[e.Caller] = v // register the vote
						voteCounter++       // increase the vote counter
						// if all players have cast a vote
						if voteCounter == players {
							// add up the votes
							var r int8 = 0
							for _, v := range votes {
								switch v {
								case Ja:
									r++
								case Nein:
									r--
								}
							}
							// if r is greater than 0 the election has passed
							if r > 0 {
								// update the term limits for the next election
								oldGov.Clear()
								oldGov.AddAll(president, chancellor)
								// todo add win condition for hitler
								state = presidentTurn // next step is to let the president choose a card to discard
								policyChoice = deck.draw(3)
								// send a successful election result and notify the cards the president has to choose from
								// in the field 'Hand'
								out <- Ok{Info: ElectionResult{
									Result: true,
									State: &GameState{
										//ElectionTracker: eTracker,
										//FascistTracker:  fTracker,
										//LiberalTracker:  lTracker,
										President:  president,
										Chancellor: chancellor,
										Hand:       append([]Policy{}, policyChoice...), // clone the policy choice
										//Roles:           append([]Role{}, roles...), // clone the roles
									},
								}}
							} else {
								state = chancellorCandidacy // next step is to start a new round
								// set the next president in line
								president = nextPresident
								// calculate the next president in a circular fashion
								nextPresident = (president + 1) % players
								// advance the election tracker
								// if advancing it triggers a forced policy enaction, do that first
								if eTracker == 2 {
									eTracker = 0
									policyChoice = deck.draw(1) // draw the policy to force
									switch policyChoice[0] {
									case LiberalPolicy:
										lTracker++
									case FascistPolicy:
										fTracker++
									} // todo handle win conditions for board trackers

									// send a failed election result and notify there was a forced policy enaction with the field
									// Hand containing the policy that was enacted
									out <- Ok{Info: ElectionResult{
										Result: false, // election failed
										State: &GameState{
											ElectionTracker: eTracker,
											FascistTracker:  fTracker,
											LiberalTracker:  lTracker,
											President:       president,
											//Chancellor: chancellor,
											Hand: append([]Policy{}, policyChoice...), // clone the policy choice
											//Roles:           append([]Role{}, roles...), // clone the roles
										},
									}}
								} else {
									eTracker++
									// send a failed election result and notify there was NOT a forced policy enaction by leaving
									// Hand nil
									out <- Ok{Info: ElectionResult{
										Result: false, // election failed
										State: &GameState{
											ElectionTracker: eTracker,
											FascistTracker:  fTracker,
											LiberalTracker:  lTracker,
											President:       president,
											//Chancellor: chancellor,
											//Hand:       append([]Policy{}, policyChoice...), // clone the policy choice
											//Roles:           append([]Role{}, roles...), // clone the roles
										},
									}}
								}
							}
						} else {
							out <- Ok{Info: General{}} // vote has been registered
						}
					} else {
						// unauthorized vote as user has already voted
						out <- Error{Err: Unauthorized{}} // send out error
					}
				} else {
					out <- Error{Err: Invalid{}} // invalid vote error
				}
			} else {
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case PolicyDiscard:
			e := event.(PolicyDiscard)
			switch state {
			case presidentTurn:
				if e.Caller == president {
					if s := e.Selection; s < 3 {
						policyChoice = append(policyChoice[:s], policyChoice[s+1:]...)
						// send a successful result and notify the chancellor has to choose from
						// the field 'Hand'
						out <- Ok{Info: PolicyDiscardOk(&GameState{
							ElectionTracker: eTracker,
							FascistTracker:  fTracker,
							LiberalTracker:  lTracker,
							President:       president,
							Chancellor:      chancellor,
							Hand:            append([]Policy{}, policyChoice...), // clone the policy choice
							//Roles:           append([]Role{}, roles...), // clone the roles
						})}
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			case chancellorTurn:
				if e.Caller == chancellor {
					if s := e.Selection; s < 2 {
						policyChoice = append(policyChoice[:s], policyChoice[s+1:]...)
						switch policyChoice[0] {
						case LiberalPolicy:
							lTracker++
						case FascistPolicy:
							fTracker++
						} // todo handle win conditions for board trackers
						// send a successful result and notify the chancellor has to choose from
						// the field 'Hand'
						out <- Ok{Info: PolicyDiscardOk(&GameState{
							ElectionTracker: eTracker,
							FascistTracker:  fTracker,
							LiberalTracker:  lTracker,
							President:       president,
							Chancellor:      chancellor,
							Hand:            append([]Policy{}, policyChoice...), // clone the policy choice
							//Roles:           append([]Role{}, roles...), // clone the roles
						})}
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			default:
				out <- Error{Err: WrongPhase{}} // send out error
			}
		}
	}
}
