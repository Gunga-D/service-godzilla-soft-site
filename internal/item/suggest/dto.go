package suggest

import "github.com/Gunga-D/service-godzilla-soft-site/internal/item"

type Suggested struct {
	Type        string
	Banner      *SuggestedBanner
	Item        *item.ItemCache
	Probability float64
}

type SuggestedBanner struct {
	Image       string
	Title       string
	Description string
	URL         string
}
