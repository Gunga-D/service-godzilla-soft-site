package algorithms

import "github.com/Gunga-D/service-godzilla-soft-site/internal/roulette"

type ItemInterval struct {
	Item  roulette.Item
	Begin float64
	End   float64
}

func NewItemInterval(item roulette.Item) ItemInterval {
	return ItemInterval{
		Item:  item,
		Begin: 0,
		End:   0,
	}
}
