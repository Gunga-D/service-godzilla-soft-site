package search

import "github.com/Gunga-D/service-godzilla-soft-site/internal/item"

type ThinkResult struct {
	HasErr     bool
	Reflection string
	Items      []item.ItemCache
}
