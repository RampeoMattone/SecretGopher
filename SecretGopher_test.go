package SecretGopher

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestNewGame(t *testing.T) {
	G := NewGame()

	if G.out == nil {
		t.Error("Output channel of game is nil")
	}

	if G.in == nil {
		t.Error("Input channel of game is nil")
	}
}

func TestHandling(t *testing.T) {
	G := NewGame()
	// Adds 10 player
	for i := 0; i < 10; i++ {
		G.in <- AddPlayer{}
		o := <-G.out
		switch o.(type) {
		case Ok:
			info := o.(Ok).Info
			switch info.(type) {
			case PlayerRegistered:
				info := info.(PlayerRegistered)
				if i == int(info) {
					t.Logf("Got PlayerRegistered for player: %d", info)
				} else {
					t.Errorf("Got wrong player number. Expected %d, got %d", i, info)
				}
			default:
				t.Error("Got wrong Info")
			}
		case Error:
			t.Error("Got wrong event")
		}
	}

	// Adds an 11th player
	G.in <- AddPlayer{}
	o := <-G.out
	switch o.(type) {
	case Ok:
		t.Error("Did not get lobby full on 11th player registered")
	case Error:
		err := o.(Error).Err
		switch err.(type) {
		case GameFull:
			t.Log("Got expected error - GameFull")
		default:
			t.Error("Got wrong Error")
		}
	}

	// Start the game
	G.in <- Start{}
	o = <-G.out
	var p, c int8
	switch o.(type) {
	case Ok:
		info := o.(Ok).Info
		switch info.(type) {
		case GameStart:
			info := info.(GameStart)
			p = info.President
			t.Log("Got expected - GameStart", info)
			// makes hitler chancellor
			/*for i, role := range info.Roles {
				if role == Hitler { c = int8(i) }
			}*/
		default:
			t.Error("Got wrong Info")
		}
	case Error:
		t.Error("Got Error - expected Gamestart")
	}

	// select the first chancellor (we select the player after the president)
	c = (p + 1) % 10
	G.in <- MakeChancellor{Caller: p, Proposal: c}
	o = <-G.out
	switch o.(type) {
	case Ok:
		info := o.(Ok).Info
		switch info.(type) {
		case ElectionStart:
			info := info.(ElectionStart)
			if c == info.Chancellor {
				t.Log("Got expected - ElectionStart", info)
			} else {
				t.Error("Got wrong Chancellor - expected", c, "got", info.Chancellor)
			}
		default:
			t.Error("Got wrong Info")
		}
	case Error:
		t.Error("Got Error - expected ElectionStart")
	}

	// send 9 out of 10 votes for yes to the gov.
	for i := 0; i < 9; i++ {
		G.in <- GovernmentVote{Caller: int8(i), Vote: Ja}
		o = <-G.out
		switch o.(type) {
		case Ok:
			info := o.(Ok).Info
			switch info.(type) {
			case VoteRegistered:
				t.Logf("Got VoteRegistered on player: %d", i)
			default:
				t.Error("Got wrong Info")
			}
		case Error:
			t.Error("Got wrong event")
		}
	}

	//send the last vote to elect the gov.
	G.in <- GovernmentVote{Caller: 9, Vote: Ja}
	o = <-G.out
	switch o.(type) {
	case Ok:
		info := o.(Ok).Info
		switch info.(type) {
		case LegislationPresident:
			info := info.(LegislationPresident)
			t.Log("Got LegislationPresident on last player", info)
			if info.State.Roles[c] == Hitler {
				t.Error("Hitler is Chancellor. game should have ended")
			}
		case GameEnd:
			t.Logf("Got GameEnd on last player")
			if info.(GameEnd).State.Roles[c] == Hitler {
				t.Logf("Hitler is Chancellor. game has succesfully ended")
				return
			}
		default:
			t.Error("Got wrong Info")
		}
	case Error:
		t.Error("Got error")
	}

	// test the policy discard system for the president (we discard the middle card)
	G.in <- PolicyDiscard{Caller: p, Selection: 1}
	o = <-G.out
	switch o.(type) {
	case Ok:
		info := o.(Ok).Info
		switch info.(type) {
		case LegislationChancellor:
			info := info.(LegislationChancellor)
			t.Log("Got expected - LegislationChancellor\n", info)
		default:
			t.Error("Got wrong Info")
		}
	case Error:
		t.Error("Got error")
	}
	// test the policy discard system for the chancellor (we discard the first card)
	G.in <- PolicyDiscard{Caller: c, Selection: 0}
	o = <-G.out
	var enacted Policy
	switch o.(type) {
	case Ok:
		info := o.(Ok).Info
		switch info.(type) {
		case PolicyEnaction:
			info := info.(PolicyEnaction)
			enacted = info.Enacted
			t.Log("Got expected - PolicyEnaction\n", info)
		default:
			t.Error("Got wrong Info", info)
		}
	case Error:
		t.Error("Got error")
	}
	// test the special power system
	if enacted == FascistPolicy {
		G.in <- SpecialPower{Caller: p, Selection: 0, Power: Investigate}
		o = <-G.out
		switch o.(type) {
		case Ok:
			info := o.(Ok).Info
			switch info.(type) {
			case SpecialPowerFeedback:
				info := info.(SpecialPowerFeedback)
				t.Log("Got expected - SpecialPowerFeedback\n", info)
			default:
				t.Error("Got wrong Info", info)
			}
		case Error:
			t.Error("Got error")
		}
	}
	// test term limits for government
	c, p = p, c
	G.in <- MakeChancellor{Caller: p, Proposal: c}
	o = <-G.out
	switch o.(type) {
	case Ok:
		t.Error("Got Ok - expected error: Invalid")
	case Error:
		err := o.(Error).Err
		switch err.(type) {
		case Invalid:
			t.Log("Got Expected error - Invalid")
		default:
			t.Error("Got wrong Error")
		}
	}
}
