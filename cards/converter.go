package cards

var combinationNames = map[int]string{
	highCard:      "highCard",
	onePair:       "onePair",
	twoPair:       "twoPair",
	threeOfAKind:  "threeOfAKind",
	straight:      "straight",
	flush:         "flush",
	fullHouse:     "fullHouse",
	fourOfAKind:   "fourOfAKind",
	straightFlush: "straightFlush",
	royalFlush:    "royalFlush",
}

var cardsNames = map[int]string{
	card2 | suitH: "2h", card2 | suitC: "2c", card2 | suitD: "2d", card2 | suitS: "2s", // 2h,2c,2d,2s
	card3 | suitH: "3h", card3 | suitC: "3c", card3 | suitD: "3d", card3 | suitS: "3s", // 3h,3c,3d,3s
	card4 | suitH: "4h", card4 | suitC: "4c", card4 | suitD: "4d", card4 | suitS: "4s", // 4h,4c,4d,4s
	card5 | suitH: "5h", card5 | suitC: "5c", card5 | suitD: "5d", card5 | suitS: "5s", // 5h,5c,5d,5s
	card6 | suitH: "6h", card6 | suitC: "6c", card6 | suitD: "6d", card6 | suitS: "6s", // 6h,6c,6d,6s
	card7 | suitH: "7h", card7 | suitC: "7c", card7 | suitD: "7d", card7 | suitS: "7s", // 7h,7c,7d,7s
	card8 | suitH: "8h", card8 | suitC: "8c", card8 | suitD: "8d", card8 | suitS: "8s", // 8h,8c,8d,8s
	card9 | suitH: "9h", card9 | suitC: "9c", card9 | suitD: "9d", card9 | suitS: "9s", // 9h,9c,9d,9s
	cardT | suitH: "Th", cardT | suitC: "Tc", cardT | suitD: "Td", cardT | suitS: "Ts", // Th,Tc,Td,Ts
	cardJ | suitH: "Jh", cardJ | suitC: "Jc", cardJ | suitD: "Jd", cardJ | suitS: "Js", // Jh,Jc,Jd,Js
	cardQ | suitH: "Qh", cardQ | suitC: "Qc", cardQ | suitD: "Qd", cardQ | suitS: "Qs", // Qh,Qc,Qd,Qs
	cardK | suitH: "Kh", cardK | suitC: "Kc", cardK | suitD: "Kd", cardK | suitS: "Ks", // Kh,Kc,Kd,Ks
	cardA | suitH: "Ah", cardA | suitC: "Ac", cardA | suitD: "Ad", cardA | suitS: "As", // Ah,Ac,Ad,As
}

var (
	cardsBinary        = map[string]int{}
	combinationsBinary = map[string]int{}
)

func init() {
	for b, n := range cardsNames {
		cardsBinary[n] = b
	}

	for b, n := range combinationNames {
		combinationsBinary[n] = b
	}
}

// GetCombinationName return combination name
func GetCombinationName(id int) string {
	return combinationNames[id]
}

// GetCardName return combination name
func GetCardName(id int) string {
	return cardsNames[id]
}

// GetCombinationID return combination id
func GetCombinationID(name string) int {
	return combinationsBinary[name]
}

// GetCardID return combination id
func GetCardID(name string) int {
	return cardsBinary[name]
}

// GetCardsNames return cards names by id
func GetCardsNames(ids []int) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = GetCardName(id)
	}

	return result
}

// GetCardsIDs return cards ids by names
func GetCardsIDs(names []string) []int {
	result := make([]int, len(names))
	for i, name := range names {
		result[i] = GetCardID(name)
	}

	return result
}
