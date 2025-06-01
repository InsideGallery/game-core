package cards

import (
	"sort"

	"github.com/InsideGallery/core/memory/set"
)

const (
	suitH = 1
	suitC = 2
	suitD = 4
	suitS = 8
	suit  = suitS | suitD | suitC | suitH
	card2 = 16
	card3 = 32
	card4 = 64
	card5 = 128
	card6 = 256
	card7 = 512
	card8 = 1024
	card9 = 2048
	cardT = 4096
	cardJ = 8192
	cardQ = 16384
	cardK = 32768
	cardA = 65536
	rank  = cardA | cardK | cardQ | cardJ | cardT | card9 | card8 | card7 | card6 | card5 | card4 | card3 | card2
)

const (
	highCard      = 131072
	onePair       = 262144
	twoPair       = 524288
	threeOfAKind  = 1048576
	straight      = 2097152
	flush         = 4194304
	fullHouse     = 8388608
	fourOfAKind   = 16777216
	straightFlush = 33554432
	royalFlush    = 67108864
)

// Combination describe combination of user
type Combination struct {
	Combination   int
	Weight        int
	KickersWeight int
	Cards         []int // always 5 cards
	Kickers       []int // from 0 to 4
}

// Score calculate score of combination
func (c *Combination) Score() int {
	return c.Combination | c.Weight | c.KickersWeight
}

// CalculateKicker calculate combination kicker
func (c *Combination) CalculateKicker(cardSet set.GenericDataSet[int]) {
	if c.Combination <= 0 {
		return
	}

	m := 5 - len(c.Cards) //nolint:mnd

	for _, c := range c.Cards {
		cardSet.Delete(c)
	}

	setCards := cardSet.ToSlice()
	if m <= 0 || len(setCards) == 0 {
		return
	}

	if m > len(setCards) {
		m = len(setCards)
	}

	sort.Slice(setCards, func(i, j int) bool {
		return setCards[i] < setCards[j]
	})

	c.Kickers = setCards[len(setCards)-m:]
	c.KickersWeight = cardsWeight(c.Kickers)
}

func cardsWeight(res []int) int {
	var weight int
	for _, c := range res {
		weight |= c & rank
	}

	return weight
}

// BinaryEvaluation evaluation by binary cards
type BinaryEvaluation struct{}

// Execute execute evaluation
func (b BinaryEvaluation) Execute(cards []int) *Combination {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i] < cards[j]
	})

	combination := b.straightFlush(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.fourOfAKind(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.fullHouse(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.flush(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.straight(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.threeOfAKind(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.twoPair(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.onePair(cards)
	if combination.Combination != 0 {
		return combination
	}

	combination = b.highCard(cards)
	if combination.Combination != 0 {
		return combination
	}

	return combination
}

func (b BinaryEvaluation) straightFlush(cards []int) *Combination {
	ranks := map[int]map[int][]int{}
	var lastRank int
	id := 0

	for _, c := range cards {
		if lastRank == c&rank {
			continue
		}

		if lastRank != 0 && lastRank<<1 != c&rank {
			id++
		}

		if _, e := ranks[id]; !e {
			ranks[id] = map[int][]int{}
		}

		ranks[id][c&suit] = append(ranks[id][c&suit], c)
		lastRank = c & rank
	}

	ace := cards[len(cards)-1]&cardA != 0
	combination := &Combination{}

	for _, suits := range ranks {
		for _, res := range suits {
			if len(res) > 5 { //nolint:mnd
				res = res[len(res)-5:]
			}

			if ace && len(res) == 4 && res[0]&card2 != 0 { //nolint:mnd
				res = append(res, cards[len(cards)-1])

				weight := cardsWeight(res)
				if weight > combination.Weight {
					combination.Combination = straightFlush
					combination.Weight = weight
					combination.KickersWeight = 0
					combination.Kickers = []int{}
					combination.Cards = res
				}
			} else if len(res) == 5 { //nolint:mnd
				weight := cardsWeight(res)
				if weight > combination.Weight {
					if res[len(res)-1]&cardA != 0 {
						combination.Combination = royalFlush
					} else {
						combination.Combination = straightFlush
					}
					combination.Weight = weight
					combination.KickersWeight = 0
					combination.Kickers = []int{}
					combination.Cards = res
				}
			}
		}
	}

	return combination
}

func (b BinaryEvaluation) fourOfAKind(cards []int) *Combination {
	ranks := map[int][]int{}
	cardSet := set.NewGenericDataSet[int]()

	for _, c := range cards {
		r := c & rank
		ranks[r] = append(ranks[r], c)
		cardSet.Add(c)
	}

	combination := &Combination{}

	for _, res := range ranks {
		if len(res) == 4 { //nolint:mnd
			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = fourOfAKind
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		}
	}

	combination.CalculateKicker(cardSet)

	return combination
}

func (b BinaryEvaluation) fullHouse(cards []int) *Combination {
	ranks := map[int][]int{}
	cardSet := set.NewGenericDataSet[int]()

	for _, c := range cards {
		r := c & rank
		ranks[r] = append(ranks[r], c)
		cardSet.Add(c)
	}

	var two [][]int
	var three [][]int

	for _, res := range ranks {
		if len(res) == 3 { //nolint:mnd
			three = append(three, res)
		}

		if len(res) == 2 { //nolint:mnd
			two = append(two, res)
		}
	}

	var mx int
	var fThree []int
	var fTwo []int

	for _, t := range three {
		if cardsWeight(t) > mx {
			mx = cardsWeight(t)
			fThree = t
		}
	}

	for _, t := range two {
		if cardsWeight(t) > mx {
			mx = cardsWeight(t)
			fTwo = t
		}
	}

	combination := &Combination{}

	if len(three) > 1 {
		three = [][]int{
			fThree,
		}
	}

	if len(two) > 1 {
		two = [][]int{
			fTwo,
		}
	}

	if len(three)+len(two) == 2 { //nolint:mnd
		var weight int
		combination.Combination = fullHouse

		for _, res := range append(three, two...) {
			weight |= cardsWeight(res)
			combination.Weight += weight
			combination.KickersWeight = 0
			combination.Kickers = []int{}
			combination.Cards = append(combination.Cards, res...)
		}
	}

	combination.CalculateKicker(cardSet)

	return combination
}

func (b BinaryEvaluation) flush(cards []int) *Combination {
	suits := map[int][]int{}
	for _, c := range cards {
		suits[c&suit] = append(suits[c&suit], c)
	}

	combination := &Combination{}

	for _, res := range suits {
		if len(res) > 5 { //nolint:mnd
			res = res[len(res)-5:]
		}

		if len(res) == 5 { //nolint:mnd
			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = flush
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		}
	}

	return combination
}

func (b BinaryEvaluation) straight(cards []int) *Combination {
	ranks := map[int][]int{}

	var lastRank int
	id := 0

	for _, c := range cards {
		if lastRank == c&rank {
			continue
		}

		if lastRank != 0 && lastRank<<1 != c&rank {
			id++
		}

		ranks[id] = append(ranks[id], c)
		lastRank = c & rank
	}

	ace := cards[len(cards)-1]&cardA != 0
	combination := &Combination{}

	for _, res := range ranks {
		if len(res) > 5 { //nolint:mnd
			res = res[len(res)-5:] //nolint:mnd
		}

		if ace && len(res) == 4 && res[0]&card2 != 0 { //nolint:mnd
			res = append(res, cards[len(cards)-1])

			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = straight
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		} else if len(res) == 5 { //nolint:mnd
			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = straight
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		}
	}

	return combination
}

func (b BinaryEvaluation) threeOfAKind(cards []int) *Combination {
	ranks := map[int][]int{}
	combination := &Combination{}
	cardSet := set.NewGenericDataSet[int]()

	for _, c := range cards {
		r := c & rank
		ranks[r] = append(ranks[r], c)
		cardSet.Add(c)
	}

	for _, res := range ranks { //nolint:mnd
		if len(res) == 3 { //nolint:mnd
			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = threeOfAKind
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		}
	}

	combination.CalculateKicker(cardSet)

	return combination
}

func (b BinaryEvaluation) twoPair(cards []int) *Combination {
	ranks := map[int][]int{}
	cardSet := set.NewGenericDataSet[int]()
	var data [][]int

	for _, c := range cards {
		r := c & rank
		ranks[r] = append(ranks[r], c)
		cardSet.Add(c)
	}

	for _, res := range ranks {
		if len(res) == 2 { //nolint:mnd
			data = append(data, res)
		}
	}

	var mx int
	var index int
	var fData []int
	var fData2 []int

	for i, t := range data {
		if cardsWeight(t) > mx {
			mx = cardsWeight(t)
			fData = t
			index = i
		}
	}

	mx = 0
	for i, t := range data {
		if i != index && cardsWeight(t) > mx {
			mx = cardsWeight(t)
			fData2 = t
		}
	}

	combination := &Combination{}

	if len(data) > 2 { //nolint:mnd
		data = [][]int{
			fData,
			fData2,
		}
	}

	if len(data) == 2 { //nolint:mnd
		var weight int
		combination.Combination = twoPair

		for _, res := range data {
			weight |= cardsWeight(res)
			combination.Weight += weight
			combination.KickersWeight = 0
			combination.Kickers = []int{}
			combination.Cards = append(combination.Cards, res...)
		}
	}

	combination.CalculateKicker(cardSet)

	return combination
}

func (b BinaryEvaluation) onePair(cards []int) *Combination {
	ranks := map[int][]int{}
	cardSet := set.NewGenericDataSet[int]()
	combination := &Combination{}

	for _, c := range cards {
		r := c & rank
		ranks[r] = append(ranks[r], c)
		cardSet.Add(c)
	}

	for _, res := range ranks {
		if len(res) == 2 { //nolint:mnd
			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = onePair
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		}
	}

	combination.CalculateKicker(cardSet)

	return combination
}

func (b BinaryEvaluation) highCard(cards []int) *Combination {
	ranks := map[int][]int{}
	cardSet := set.NewGenericDataSet[int]()

	for _, c := range cards {
		r := c & rank
		ranks[r] = append(ranks[r], c)
		cardSet.Add(c)
	}

	combination := &Combination{}

	for _, res := range ranks {
		if len(res) == 1 {
			weight := cardsWeight(res)
			if weight > combination.Weight {
				combination.Combination = highCard
				combination.Weight = weight
				combination.KickersWeight = 0
				combination.Kickers = []int{}
				combination.Cards = res
			}
		}
	}

	combination.CalculateKicker(cardSet)

	return combination
}
