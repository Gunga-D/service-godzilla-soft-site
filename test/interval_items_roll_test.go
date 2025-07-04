package test

import (
	"fmt"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/algorithms"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/roulette"
	"testing"
)

func GetItems() []roulette.Item {
	var items = make([]roulette.Item, 20)
	for indx, _ := range items {
		items[indx] = roulette.Item{
			Id:           0,
			TotalCost:    10 + int64(indx)*10,
			ItemCategory: roulette.Golden,
		}
	}
	return items
}
func TestRollItems(t *testing.T) {
	items := GetItems()
	arr := algorithms.NewIntervalArray(items, 100)
	arr.Split()

	probs := arr.EstimateProbabilities()
	totalChance := float64(0)
	for i := 0; i < len(probs); i++ {
		totalChance += probs[i]
	}
	fmt.Println(probs)
}
