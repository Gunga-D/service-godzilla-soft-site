package deepthink

import "github.com/Gunga-D/service-godzilla-soft-site/internal/item"

type ThinkResult struct {
	Reflection string
	Items      []item.ItemCache
}
