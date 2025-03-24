package search_suggest

type SearchSuggestRequest struct {
	Query string `json:"query"`
}

type SearchSuggestResponsePayload struct {
	Items []SearchSuggestDTO `json:"items"`
}

type SearchSuggestDTO struct {
	ItemID           int64    `json:"item_id"`
	ItemCategoryID   int64    `json:"item_category_id"`
	ItemTitle        string   `json:"item_title"`
	ItemCurrentPrice float64  `json:"item_current_price"`
	ItemIsForSale    bool     `json:"item_is_for_sale"`
	ItemOldPrice     *float64 `json:"item_old_price"`
	ItemThumbnailURL string   `json:"item_thumbnail_url"`
	ItemType         string   `json:"item_type"`
	Probability      float64  `json:"probability"`
}
