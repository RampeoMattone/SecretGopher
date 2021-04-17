package SecretGopher

import "math/rand"

type deck struct {
	d [17]Policy // d is the stack of cards
	p uint8      // p is the position within the stack
}

// newDeck generates a new deck for a game and shuffles it ahead of time
func newDeck() deck {
	var d = deck{
		p: 0,
		d: [17]Policy{
			LiberalPolicy, LiberalPolicy, LiberalPolicy, LiberalPolicy, LiberalPolicy, LiberalPolicy,
			FascistPolicy, FascistPolicy, FascistPolicy, FascistPolicy, FascistPolicy, FascistPolicy,
			FascistPolicy, FascistPolicy, FascistPolicy, FascistPolicy, FascistPolicy,
		},
	}
	d.shuffle()
	return d
}

// shuffle shuffles the elements of the deck pseudorandomically
func (d *deck) shuffle() {
	d.p = 0
	rand.Shuffle(17, func(i, j int) {
		d.d[i], d.d[j] = d.d[j], d.d[i]
	})
}

// draw draws the top n cards from the deck, making sure to move them away and not draw them again
func (d *deck) draw(n uint8) []Policy {
	var r = make([]Policy, n)
	for i := uint8(0); i < n; i++ {
		r[i] = d.d[d.p]
		d.p++
	}
	// "If there are fewer than three tiles remaining in the policy deck at the end of a Legislative Session,
	// they are shuffled with the Discard pile to create a new policy deck. Unused policy tiles are not revealed."
	if d.p > 14 {
		d.shuffle()
	}
	return r
}

// peek reveals the top 3 cards from the deck, making sure to leave the deck unaltered
func (d *deck) peek() [3]Policy {
	var r = [3]Policy{}
	for i := uint8(0); i < 3; i++ {
		r[i] = d.d[d.p+i]
	}
	return r
}
