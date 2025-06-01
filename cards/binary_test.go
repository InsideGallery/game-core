package cards

import (
	"log/slog"
	"testing"

	"github.com/InsideGallery/core/testutils"
)

func TestExecution(t *testing.T) {
	b := BinaryEvaluation{}
	comb := b.Execute([]int{
		card4 | suitH,
		card4 | suitC,
		card5 | suitD,
		card5 | suitS,
		card5 | suitC,
		cardJ | suitC,
		cardA | suitD,
		cardA | suitC,
	})
	testutils.Equal(t, comb.Combination, fullHouse)
	comb2 := b.Execute([]int{
		card4 | suitH,
		card4 | suitC,
		card4 | suitD,
		card5 | suitS,
		card5 | suitC,
		cardJ | suitC,
		cardA | suitD,
		cardA | suitC,
	})
	testutils.Equal(t, comb2.Combination, fullHouse)
	slog.Default().Error("Scores", "comb2.Score()", comb2.Score(), "comb.Score()", comb.Score())
	testutils.Equal(t, comb2.Score() < comb.Score(), true)
}

func TestCornerCases(t *testing.T) {
	testcases := []struct {
		name         string
		cards1       []string
		cards2       []string
		win1         bool
		split        bool
		combination1 string
		combination2 string
	}{
		{
			name:         "straight vs onePair",
			cards1:       []string{"6d", "3c", "2h", "3d", "4c", "5s", "8d"},
			cards2:       []string{"7d", "4c", "2h", "3d", "4c", "5s", "8d"},
			win1:         true,
			split:        false,
			combination1: "straight",
			combination2: "onePair",
		},
		{
			name:         "straight vs straight split",
			cards1:       []string{"6d", "3c", "2h", "3d", "4c", "5s", "8d"},
			cards2:       []string{"6s", "4c", "2h", "3d", "4c", "5s", "8d"},
			win1:         false,
			split:        true,
			combination1: "straight",
			combination2: "straight",
		},
		{
			name:         "straight vs straight",
			cards1:       []string{"2d", "3c", "4h", "5d", "6c", "7s", "2s"},
			cards2:       []string{"2s", "8c", "4h", "5d", "6c", "7s", "2s"},
			win1:         false,
			split:        false,
			combination1: "straight",
			combination2: "straight",
		},
		{
			name:         "straight vs straight",
			cards1:       []string{"2d", "3c", "4h", "7d", "8c", "3s", "2s"},
			cards2:       []string{"3h", "3d", "4h", "7d", "8c", "3s", "2s"},
			win1:         false,
			split:        false,
			combination1: "twoPair",
			combination2: "threeOfAKind",
		},
		{
			name:         "flush vs highCard",
			cards1:       []string{"2d", "3d", "4d", "7d", "8d", "9s", "2s"},
			cards2:       []string{"3h", "Ad", "4h", "7d", "8c", "9s", "2s"},
			win1:         true,
			split:        false,
			combination1: "flush",
			combination2: "highCard",
		},
	}

	ev := BinaryEvaluation{}
	for _, test := range testcases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			comb1 := ev.Execute(GetCardsIDs(test.cards1))
			comb2 := ev.Execute(GetCardsIDs(test.cards2))
			testutils.Equal(t, GetCombinationName(comb1.Combination), test.combination1)
			testutils.Equal(t, GetCombinationName(comb2.Combination), test.combination2)
			testutils.Equal(t, comb1.Score() > comb2.Score(), test.win1)
			testutils.Equal(t, comb1.Score() == comb2.Score(), test.split)
		})
	}
}

/*
BenchmarkExecution-12    	  238504	      5023 ns/op
*/

func BenchmarkExecution(b *testing.B) {
	ev := BinaryEvaluation{}
	cards := GetCardsIDs([]string{"7d", "4c", "2h", "3d", "9c", "5s", "8d"})
	for i := 0; i < b.N; i++ {
		ev.Execute(cards)
	}
}
