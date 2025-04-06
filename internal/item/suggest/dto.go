package suggest

type Suggested struct {
	Type        string
	Banner      *SuggestedBanner
	Item        *SuggestedItem
	Probability float64
}

type SuggestedBanner struct {
	Image string
	Title string
	URL   string
}

type SuggestedItem struct {
	ID           int64
	CategoryID   int64
	Title        string
	CurrentPrice int64
	IsForSale    bool
	OldPrice     *int64
	ThumbnailURL string
	IsSteamGift  bool
}
