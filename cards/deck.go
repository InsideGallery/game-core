package cards

import (
	"errors"
	"math/rand/v2"
	"sync"
)

var ErrNoCardsInDeck = errors.New("no cards in deck")

var deck = []int{
	card2 | suitH, card2 | suitC, card2 | suitD, card2 | suitS, // 2h,2c,2d,2s
	card3 | suitH, card3 | suitC, card3 | suitD, card3 | suitS, // 3h,3c,3d,3s
	card4 | suitH, card4 | suitC, card4 | suitD, card4 | suitS, // 4h,4c,4d,4s
	card5 | suitH, card5 | suitC, card5 | suitD, card5 | suitS, // 5h,5c,5d,5s
	card6 | suitH, card6 | suitC, card6 | suitD, card6 | suitS, // 6h,6c,6d,6s
	card7 | suitH, card7 | suitC, card7 | suitD, card7 | suitS, // 7h,7c,7d,7s
	card8 | suitH, card8 | suitC, card8 | suitD, card8 | suitS, // 8h,8c,8d,8s
	card9 | suitH, card9 | suitC, card9 | suitD, card9 | suitS, // 9h,9c,9d,9s
	cardT | suitH, cardT | suitC, cardT | suitD, cardT | suitS, // Th,Tc,Td,Ts
	cardJ | suitH, cardJ | suitC, cardJ | suitD, cardJ | suitS, // Jh,Jc,Jd,Js
	cardQ | suitH, cardQ | suitC, cardQ | suitD, cardQ | suitS, // Qh,Qc,Qd,Qs
	cardK | suitH, cardK | suitC, cardK | suitD, cardK | suitS, // Kh,Kc,Kd,Ks
	cardA | suitH, cardA | suitC, cardA | suitD, cardA | suitS, // Ah,Ac,Ad,As
}

// Deck contains cards from current game
type Deck struct {
	cards []int
	mu    sync.Mutex
}

// NewDeck return new deck
func NewDeck() *Deck {
	d := &Deck{
		cards: make([]int, len(deck)),
	}
	copy(d.cards, deck)

	return d
}

// GetRandomCard return random card, and remove from deck
func (d *Deck) GetRandomCard() (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.cards) == 0 {
		return 0, ErrNoCardsInDeck
	}

	i := rand.IntN(len(d.cards)) // nolint:gosec
	card := d.cards[i]
	copy(d.cards[i:], d.cards[i+1:])
	d.cards = d.cards[:len(d.cards)-1]

	return card, nil
}

// Shuffle shuffle deck
func (d *Deck) Shuffle() {
	d.mu.Lock()
	defer d.mu.Unlock()

	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

// Top return card from top
func (d *Deck) Top() (card int, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.cards) == 0 {
		return 0, ErrNoCardsInDeck
	}

	card, d.cards = d.cards[0], d.cards[1:]

	return card, nil
}

// GetTopCards return top cards, and remove from deck
func (d *Deck) GetTopCards(n int) ([]int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if n == -1 {
		n = len(d.cards)
	}

	if n <= 0 {
		return []int{}, nil
	}

	result := make([]int, n)
	copy(result, d.cards[:n])

	return result, nil
}
