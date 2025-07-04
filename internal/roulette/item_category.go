package roulette

type ItemCategory = string

const (
	Common   ItemCategory = "common"
	Uncommon ItemCategory = "uncommon"
	Rare     ItemCategory = "rare"
	Special  ItemCategory = "special"
	Golden   ItemCategory = "golden"
)

func CategoryFromPrice(price int64) ItemCategory {
	if price <= 10000 {
		return Common
	} else if price <= 25000 {
		return Uncommon
	} else if price <= 50000 {
		return Rare
	} else if price <= 100000 {
		return Special
	}
	return Golden
}
