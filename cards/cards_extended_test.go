package cards

import (
	"testing"

	"github.com/InsideGallery/core/testutils"
)

// ---------- Deck tests ----------

func TestNewDeck(t *testing.T) {
	d := NewDeck()
	testutils.Equal(t, len(d.cards), 52)
}

func TestDeckShuffle(t *testing.T) {
	d := NewDeck()
	original := make([]int, len(d.cards))
	copy(original, d.cards)
	d.Shuffle()
	// After shuffle the deck still has 52 cards
	testutils.Equal(t, len(d.cards), 52)
	// Very unlikely that the order is exactly the same after shuffling
	// but we just check the size is right.
}

func TestDeckTop(t *testing.T) {
	d := NewDeck()
	first := d.cards[0]
	card, err := d.Top()
	testutils.Equal(t, err, nil)
	testutils.Equal(t, card, first)
	testutils.Equal(t, len(d.cards), 51)
}

func TestDeckTopEmpty(t *testing.T) {
	d := &Deck{cards: []int{}}
	_, err := d.Top()
	testutils.Equal(t, err, ErrNoCardsInDeck)
}

func TestDeckGetTopCardsNormal(t *testing.T) {
	d := NewDeck()
	cards, err := d.GetTopCards(5)
	testutils.Equal(t, err, nil)
	testutils.Equal(t, len(cards), 5)
	testutils.Equal(t, len(d.cards), 47)
}

func TestDeckGetTopCardsAll(t *testing.T) {
	d := NewDeck()
	cards, err := d.GetTopCards(-1)
	testutils.Equal(t, err, nil)
	testutils.Equal(t, len(cards), 52)
	testutils.Equal(t, len(d.cards), 0)
}

func TestDeckGetTopCardsZero(t *testing.T) {
	d := NewDeck()
	cards, err := d.GetTopCards(0)
	testutils.Equal(t, err, nil)
	testutils.Equal(t, len(cards), 0)
}

func TestDeckGetTopCardsTooMany(t *testing.T) {
	d := NewDeck()
	_, err := d.GetTopCards(53)
	testutils.Equal(t, err, ErrNoCardsInDeck)
}

func TestDeckGetRandomCard(t *testing.T) {
	d := NewDeck()
	card, err := d.GetRandomCard()
	testutils.Equal(t, err, nil)
	testutils.Equal(t, card > 0, true)
	testutils.Equal(t, len(d.cards), 51)
}

func TestDeckGetRandomCardUntilEmpty(t *testing.T) {
	d := NewDeck()
	for i := 0; i < 52; i++ {
		_, err := d.GetRandomCard()
		testutils.Equal(t, err, nil)
	}
	_, err := d.GetRandomCard()
	testutils.Equal(t, err, ErrNoCardsInDeck)
}

// ---------- Converter tests ----------

func TestGetCombinationName(t *testing.T) {
	tests := []struct {
		id   int
		name string
	}{
		{highCard, "highCard"},
		{onePair, "onePair"},
		{twoPair, "twoPair"},
		{threeOfAKind, "threeOfAKind"},
		{straight, "straight"},
		{flush, "flush"},
		{fullHouse, "fullHouse"},
		{fourOfAKind, "fourOfAKind"},
		{straightFlush, "straightFlush"},
		{royalFlush, "royalFlush"},
		{0, ""},
	}
	for _, tc := range tests {
		testutils.Equal(t, GetCombinationName(tc.id), tc.name)
	}
}

func TestGetCardName(t *testing.T) {
	testutils.Equal(t, GetCardName(card2|suitH), "2h")
	testutils.Equal(t, GetCardName(cardA|suitS), "As")
	testutils.Equal(t, GetCardName(0), "")
}

func TestGetCombinationID(t *testing.T) {
	testutils.Equal(t, GetCombinationID("highCard"), highCard)
	testutils.Equal(t, GetCombinationID("royalFlush"), royalFlush)
	testutils.Equal(t, GetCombinationID("nonexistent"), 0)
}

func TestGetCardID(t *testing.T) {
	testutils.Equal(t, GetCardID("2h"), card2|suitH)
	testutils.Equal(t, GetCardID("As"), cardA|suitS)
	testutils.Equal(t, GetCardID("nonexistent"), 0)
}

func TestGetCardsNames(t *testing.T) {
	ids := []int{card2 | suitH, cardA | suitS}
	names := GetCardsNames(ids)
	testutils.Equal(t, names, []string{"2h", "As"})
}

func TestGetCardsNamesEmpty(t *testing.T) {
	names := GetCardsNames([]int{})
	testutils.Equal(t, len(names), 0)
}

func TestGetCardsIDs(t *testing.T) {
	names := []string{"2h", "As"}
	ids := GetCardsIDs(names)
	testutils.Equal(t, ids, []int{card2 | suitH, cardA | suitS})
}

// ---------- Combination tests ----------

func TestCombinationScore(t *testing.T) {
	c := &Combination{
		Combination:   highCard,
		Weight:        cardA,
		KickersWeight: cardK,
	}
	testutils.Equal(t, c.Score(), highCard|cardA|cardK)
}

func TestCombinationScoreZero(t *testing.T) {
	c := &Combination{}
	testutils.Equal(t, c.Score(), 0)
}

func TestCalculateKickerNoCombination(t *testing.T) {
	c := &Combination{Combination: 0}
	// CalculateKicker should return immediately when Combination <= 0
	// Just verify it doesn't panic
	c.CalculateKicker(nil)
}

// ---------- Binary evaluation: straight flush, four of a kind, straight, ace-low ----------

func TestStraightFlush(t *testing.T) {
	b := BinaryEvaluation{}
	// Straight flush: 5h 6h 7h 8h 9h + filler
	comb := b.Execute(GetCardsIDs([]string{"5h", "6h", "7h", "8h", "9h", "2c", "3d"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "straightFlush")
}

func TestRoyalFlush(t *testing.T) {
	b := BinaryEvaluation{}
	// Royal flush: Th Jh Qh Kh Ah + filler
	comb := b.Execute(GetCardsIDs([]string{"Th", "Jh", "Qh", "Kh", "Ah", "2c", "3d"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "royalFlush")
}

func TestAceLowStraight(t *testing.T) {
	b := BinaryEvaluation{}
	// Ace-low straight: Ah 2c 3d 4s 5h + filler
	comb := b.Execute(GetCardsIDs([]string{"Ah", "2c", "3d", "4s", "5h", "9c", "Td"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "straight")
}

func TestAceLowStraightFlush(t *testing.T) {
	b := BinaryEvaluation{}
	// Ace-low straight flush: Ah 2h 3h 4h 5h + filler
	comb := b.Execute(GetCardsIDs([]string{"Ah", "2h", "3h", "4h", "5h", "9c", "Td"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "straightFlush")
}

func TestFourOfAKind(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute(GetCardsIDs([]string{"Ah", "Ac", "Ad", "As", "5h", "9c", "Td"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "fourOfAKind")
	testutils.Equal(t, len(comb.Cards), 4)
	testutils.Equal(t, len(comb.Kickers), 1)
}

func TestThreeOfAKind(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute(GetCardsIDs([]string{"Ah", "Ac", "Ad", "5s", "7h", "9c", "Td"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "threeOfAKind")
	testutils.Equal(t, len(comb.Cards), 3)
	testutils.Equal(t, len(comb.Kickers), 2)
}

func TestOnePair(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute(GetCardsIDs([]string{"Ah", "Ac", "3d", "5s", "7h", "9c", "Td"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "onePair")
	testutils.Equal(t, len(comb.Cards), 2)
	testutils.Equal(t, len(comb.Kickers), 3)
}

func TestHighCard(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute(GetCardsIDs([]string{"Ah", "2c", "4d", "6s", "8h", "Tc", "Qd"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "highCard")
}

func TestFlush(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute(GetCardsIDs([]string{"2h", "4h", "6h", "8h", "Th", "3c", "5d"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "flush")
}

func TestTwoPair(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute(GetCardsIDs([]string{"Ah", "Ac", "Kd", "Ks", "7h", "9c", "Td"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "twoPair")
	testutils.Equal(t, len(comb.Cards), 4)
	testutils.Equal(t, len(comb.Kickers), 1)
}

func TestStraightVsStraightWithHighCards(t *testing.T) {
	b := BinaryEvaluation{}
	// 6-high straight vs 7-high straight
	comb1 := b.Execute(GetCardsIDs([]string{"2d", "3c", "4h", "5d", "6c", "9s", "Kd"}))
	comb2 := b.Execute(GetCardsIDs([]string{"3s", "4c", "5h", "6d", "7c", "9s", "Kd"}))
	testutils.Equal(t, GetCombinationName(comb1.Combination), "straight")
	testutils.Equal(t, GetCombinationName(comb2.Combination), "straight")
	testutils.Equal(t, comb2.Score() > comb1.Score(), true)
}

func TestMoreThanFiveFlushCards(t *testing.T) {
	b := BinaryEvaluation{}
	// 6 cards of same suit - should pick top 5
	comb := b.Execute(GetCardsIDs([]string{"2h", "4h", "6h", "8h", "Th", "Qh", "3c"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "flush")
	testutils.Equal(t, len(comb.Cards), 5)
}

func TestMoreThanFiveStraightCards(t *testing.T) {
	b := BinaryEvaluation{}
	// 6 card straight
	comb := b.Execute(GetCardsIDs([]string{"2d", "3c", "4h", "5d", "6c", "7s", "Kd"}))
	testutils.Equal(t, GetCombinationName(comb.Combination), "straight")
	testutils.Equal(t, len(comb.Cards), 5)
}
