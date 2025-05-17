package category

import "github.com/AlekSi/pointer"

const (
	GamesCategoryID = 10001
)

type Category struct {
	ID       int64
	ParentID *int64
	Name     string
}

var CategoriesByID = map[int64]Category{
	10001: {
		ID:       10001,
		ParentID: nil,
		Name:     "Игры",
	},
	10002: {
		ID:       10002,
		ParentID: nil,
		Name:     "Подписки",
	},
	10003: {
		ID:       10003,
		ParentID: pointer.ToInt64(10002),
		Name:     "Xbox",
	},
	10004: {
		ID:       10004,
		ParentID: nil,
		Name:     "Пополнения",
	},
	10005: {
		ID:       10005,
		ParentID: pointer.ToInt64(10004),
		Name:     "Xbox",
	},
	10006: {
		ID:       10006,
		ParentID: pointer.ToInt64(10004),
		Name:     "Playstation",
	},
	10007: {
		ID:       10007,
		ParentID: pointer.ToInt64(10004),
		Name:     "Fortnite",
	},
	10008: {
		ID:       10008,
		ParentID: pointer.ToInt64(10004),
		Name:     "Roblox",
	},
	10009: {
		ID:       10009,
		ParentID: pointer.ToInt64(10004),
		Name:     "Valorant",
	},
	10010: {
		ID:       10010,
		ParentID: pointer.ToInt64(10004),
		Name:     "Nintendo",
	},
}
