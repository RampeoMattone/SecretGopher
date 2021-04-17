package SecretGopher

import (
	"math/rand"
)

// NewGame creates a game structure and subscribes a goroutine to listen to the events for the game
func NewGame() Game {
	in := make(chan Event)
	out := make(chan Output)
	g := gameData{}
	go g.handleGame(in, out)
	return Game{
		in:  in,
		out: out,
	}
}

type gameData struct {
	state         state
	players       int8
	deck          deck
	president     int8
	chancellor    int8
	roles         []Role
	nextPresident int8
	oldGov        set
	investigated  set
	votes         []Vote
	voted		  int8
	killed        []int8
	policyChoice  []Policy
	eTracker      int8
	fTracker      int8
	lTracker      int8
}

// search Returns a boolean value describing if the element exists in arr
func search(arr []int8, elem int8) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func (g gameData) gameOver() GameEnding {
	// check the fascist policies
	if g.fTracker == 6 {
		return FascistPolicyWin
		// check the liberal policies
	} else if g.lTracker == 5 {
		return LiberalPolicyWin
		// check if hitler is chancellor
	} else if g.roles[g.chancellor] == Hitler {
		return FascistElectionWin
		// check if hitler is dead
	} else {
		for hitler, v := range g.roles { // find the player who is hitler
			if v == Hitler {
				if search(g.killed, int8(hitler)) {
					return LiberalExecutionWin
				}
				break
			}
		}
		// if none of the above succeed then the game is still running
		return StillRunning
	}
}

func (g gameData) enactPolicyInactive() SpecialPowers {
	s := Nothing // special powers checked only when the policy is fascist
	switch g.policyChoice[0] {
	case LiberalPolicy:
		g.lTracker++
	case FascistPolicy:
		g.fTracker++
		switch g.players {
		case 5, 6:
			s = powersTable[0][g.fTracker]
		case 7, 8:
			s = powersTable[1][g.fTracker]
		case 9, 10:
			s = powersTable[2][g.fTracker]
		}
	}
	return s
}

func (g gameData) enactPolicyActive(out chan<- Output) {
	s := g.enactPolicyInactive() // s is the special power
	// checks if the game is over (if the policy limit for a party has been reached)
	if o := g.gameOver(); o != StillRunning {
		g.state = gameEnd
		out <- Ok{Info: GameEnd{
			Why: o,
			State: g.shareState(),
		}}
		return
	}
	// update the state of the game in accordance to the special power
	switch s {
	case Nothing:
		g.state = chancellorCandidacy
		// set the next president in line
		g.president = g.nextPresident
		// calculate the next president in a circular fashion
		g.nextPresident = (g.president + 1) % g.players
	case Execution:
		g.state = specialExecution
	case Election:
		g.state = specialElection
	case Investigate:
		g.state = specialInvestigate
	case Peek:
		g.state = specialPeek
	}

	// send a successful result for the enaction of a policy
	out <- Ok{Info: PolicyEnaction{
		Enacted:      g.policyChoice[0],
		SpecialPower: s,
		State: g.shareState(),
	}}
}

func (g gameData) inactiveGov(out chan<- Output) {
	g.state = chancellorCandidacy // next step is to start a new round
	// set the next president in line
	g.president = g.nextPresident
	// calculate the next president in a circular fashion
	g.nextPresident = (g.president + 1) % g.players
	// advance the election tracker
	// if advancing it triggers a forced policy enaction, do that first
	if g.eTracker == 2 {
		g.eTracker = 0
		g.policyChoice = g.deck.draw(1) // draw the policy to force
		g.enactPolicyInactive()

		// checks if the game is over (if the policy limit for a party has been reached)
		if o := g.gameOver(); o != StillRunning {
			g.state = gameEnd
			out <- Ok{Info: GameEnd{
				Why: o,
				State: g.shareState(),
			}}
			return
		}

		// send a failed election result and notify there was a forced policy enaction
		out <- Ok{Info: PolicyEnaction{
			Enacted:      g.policyChoice[0], // election failed
			SpecialPower: Nothing,
			State: g.shareState(),
		}}
	} else {
		g.eTracker++
		// send a failed election result and notify there was NOT a forced policy enaction by leaving
		// Hand nil
		out <- Ok{Info: NextPresident(g.shareState())}
	}
}

func (g gameData) shareState() GameState {
	return GameState{
		ElectionTracker: g.eTracker,
		FascistTracker:  g.fTracker,
		LiberalTracker:  g.lTracker,
		President:       g.president,
		Chancellor:      g.chancellor,
		Roles:           append([]Role{}, g.roles...), // clone the roles
		Votes:           append([]Vote{}, g.votes...), // clone the votes
		Killed:          append([]int8{}, g.killed...),
	}
}

// handleGame handles the game events.
func (g gameData) handleGame(in <-chan Event, out chan<- Output) {
	defer close(out)
	g = gameData{
		state:   waitingPlayers,
		players: 0,
		//deck:          // initialized with Start
		president:  NotSet,
		chancellor: NotSet,
		//roles:         // initialized with Start
		nextPresident: NotSet,
		oldGov:        make(set, 2),
		killed:        make([]int8, 0, 2),
		//investigated:  // initialized with Start
		//votes:         // initialized with Start
		voted: 0,
		//policyChoice:  // initialized when entering presidentLegislation
		eTracker: 0,
		fTracker: 0,
		lTracker: 0,
	}
	for {
		event := <-in
		switch event.(type) {
		case AddPlayer:
			// if the game is accepting players
			if g.state == waitingPlayers {
				if g.players < 10 {
					g.players++                                      // adds a player to the game
					out <- Ok{Info: PlayerRegistered(g.players - 1)} // say the player was registered under the player number
				} else {
					out <- Error{Err: GameFull{}} // send out error
				}
			} else {
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case Start:
			// if the game was accepting players
			if g.state == waitingPlayers {
				if g.players >= 5 {
					g.roles = make([]Role, g.players)    // initialize roles to the proper size
					g.votes = make([]Vote, g.players) // initialize votes to the proper size
					g.oldGov = make(set, 2)              // initialize oldGov to the proper size
					g.deck = newDeck()                   // initialize deck and shuffle it

					g.roles[rand.Intn(int(g.players))] = Hitler // set one player as Hitler
					var nF int                                  // number of fascists based on the lobby size
					switch g.players {
					case 5, 6:
						g.investigated = nil
						nF = 1
					case 7, 8:
						g.investigated = make(set, 1)
						nF = 2
					case 9, 10:
						g.investigated = make(set, 2)
						nF = 3
					}
					// assign nF FascistParty roles randomly
					for i := 0; i < nF; {
						// extract a player
						// if the role for that player is not FascistParty or Hitler, set him as FascistParty
						// and increase the counter
						if r := rand.Intn(int(g.players)); g.roles[r] == LiberalParty {
							g.roles[r] = FascistParty
							i++
						}
					}
					// the first player to be president is random
					g.president = int8(rand.Intn(int(g.players)))
					// set the next president in line
					g.nextPresident = (g.president + 1) % g.players

					g.state = chancellorCandidacy // after a president is selected, a chancellor needs to be selected

					out <- Ok{Info: GameStart(g.shareState())} // tell the caller the game has started
				}
			} else {
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case MakeChancellor:
			// if the game was accepting players
			if g.state == chancellorCandidacy {
				e := event.(MakeChancellor)
				if e.Caller == g.president {
					if !g.oldGov.has(e.Proposal) {
						g.chancellor = e.Proposal
						g.state = governmentElection
						g.votes = make([]Vote, g.players) // reset votes
						g.voted = 0
						out <- Ok{Info: ElectionStart(g.shareState())} // say the chancellor registration was successful
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
			if g.state == governmentElection {
				e := event.(GovernmentVote)
				// check that the vote is valid
				if v := e.Vote; v == Ja || v == Nein {
					// if the user hasn't voted yet
					if g.votes[e.Caller] == NoVote {
						g.voted++
						g.votes[e.Caller] = v // register the vote
						// if all players have cast a vote
						if g.voted == g.players {
							// add up the votes
							var r int8 = 0
							for _, v := range g.votes {
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
								g.oldGov.clear()
								g.oldGov.addAll(g.president, g.chancellor)

								// checks if the game is over (if hitler is chancellor)
								if o := g.gameOver(); o != StillRunning {
									out <- Ok{Info: GameEnd{
										Why: o,
										State: g.shareState(),
									}}
									return // end the game
								}
								g.state = presidentLegislation // next step is to let the president choose a card to discard
								g.policyChoice = g.deck.draw(3)
								// send a successful election result and notify the cards the president has to choose from
								// in the field 'Hand'
								out <- Ok{Info: LegislationPresident{
									Hand: append([]Policy{}, g.policyChoice...), // clone the policy choice
									State: g.shareState(),
								}}
							} else {
								g.inactiveGov(out) // gov was inactive, apply rules and effects
								if g.state == gameEnd {
									return // stops the handler
								}
							}
						} else {
							out <- Ok{Info: VoteRegistered{}} // vote has been registered
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
			switch g.state {
			case presidentLegislation:
				if e.Caller == g.president {
					if s := e.Selection; s < 3 {
						g.policyChoice = append(g.policyChoice[:s], g.policyChoice[s+1:]...)
						// send a successful result and notify the chancellor has to choose from
						// the field 'Hand'
						out <- Ok{Info: LegislationChancellor{
							Hand: append([]Policy{}, g.policyChoice...), // clone the policy choice
							State: g.shareState(),
						}}
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			case chancellorLegislation:
				if e.Caller == g.chancellor {
					if s := e.Selection; s < 2 {
						g.policyChoice = append(g.policyChoice[:s], g.policyChoice[s+1:]...)
						if g.fTracker == 5 {
							// send out a veto request
							g.state = vetoChancellor
							out <- Ok{Info: VetoRequest(g.shareState())}
						} else {
							s := g.enactPolicyInactive() // s is the special power
							// checks if the game is over (if the policy limit for a party has been reached)
							if o := g.gameOver(); o != StillRunning {
								out <- Ok{Info: GameEnd{
									Why: o,
									State: g.shareState(),
								}}
								return
							}
							// update the state of the game in accordance to the special power
							switch s {
							case Nothing:
								g.state = chancellorCandidacy
								// set the next president in line
								g.president = g.nextPresident
								// calculate the next president in a circular fashion
								g.nextPresident = (g.president + 1) % g.players
							case Execution:
								g.state = specialExecution
							case Election:
								g.state = specialElection
							case Investigate:
								g.state = specialInvestigate
							case Peek:
								g.state = specialPeek
							}

							// send a successful result for the enaction of a policy
							out <- Ok{Info: PolicyEnaction{
								Enacted:      g.policyChoice[0],
								SpecialPower: s,
								State: g.shareState(),
							}}
						}
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			default:
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case VetoResponse:
			e := event.(VetoResponse)
			switch g.state {
			case vetoChancellor:
				if e.Caller == g.chancellor {
					if e.Approves {
						g.state = vetoPresident
						out <- Ok{Info: VetoRequest(g.shareState())}
					} else {
						g.enactPolicyActive(out)
						if g.state == gameEnd {
							return // stops the handler
						}
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			case vetoPresident:
				if e.Caller == g.president {
					if e.Approves {
						g.inactiveGov(out) // gov was inactive, apply rules and effects
					} else {
						g.enactPolicyActive(out)
					}
					if g.state == gameEnd {
						return // stops the handler
					}
				} else {
					out <- Error{Err: Unauthorized{}} // send out error
				}
			default:
				out <- Error{Err: WrongPhase{}} // send out error
			}
		case SpecialPower:
			e := event.(SpecialPower)
			if e.Caller == g.president {
				switch e.Power {
				case Peek:
					if g.state == specialPeek {
						out <- Ok{Info: SpecialPowerFeedback{
							Feedback: g.deck.peek(),
							State: g.shareState(),
						}} // send out error
					} else {
						out <- Error{Err: WrongPhase{}} // send out error
					}
				case Election:
					if g.state == specialElection {
						// the president cannot choose himself
						if e.Selection < g.players && e.Selection != g.president {
							g.president = e.Selection
							g.state = chancellorCandidacy
							out <- Ok{Info: SpecialPowerFeedback{
								Feedback: g.deck.peek(),
								State: g.shareState(),
							}}
						} else {
							out <- Error{Err: Invalid{}} // send out error
						}
					} else {
						out <- Error{Err: WrongPhase{}} // send out error
					}
				case Execution:
					if g.state == specialExecution {
						if e.Selection < g.players && !search(g.killed, e.Selection) {
							g.killed[len(g.killed)] = e.Selection
							g.state = chancellorCandidacy
							out <- Ok{Info: SpecialPowerFeedback{
								State: g.shareState(),
							}}
						} else {
							out <- Error{Err: Invalid{}} // send out error
						}
					} else {
						out <- Error{Err: WrongPhase{}} // send out error
					}
				case Investigate:
					if g.state == specialInvestigate {
						if e.Selection < g.players && !g.investigated.has(e.Selection) {
							g.investigated.add(e.Selection)
							g.state = chancellorCandidacy
							out <- Ok{Info: SpecialPowerFeedback{
								Feedback: g.roles[e.Selection],
								State: g.shareState(),
							}}
						} else {
							out <- Error{Err: Invalid{}} // send out error
						}
					} else {
						out <- Error{Err: WrongPhase{}} // send out error
					}
				default:
					out <- Error{Err: Invalid{}} // send out error
				}
			} else {
				out <- Error{Err: Unauthorized{}} // send out error
			}
		default:
			out <- Error{Err: Invalid{}} // send out error for invalid event
		}
	}
}
