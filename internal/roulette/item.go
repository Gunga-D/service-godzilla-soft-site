package roulette

type Item struct {
	Id           int64        `json:"id" db:"id"`
	TotalCost    int64        `json:"total_cost" db:"total_cost"`
	ItemCategory ItemCategory `json:"item_category" db:"item_category"`
}
