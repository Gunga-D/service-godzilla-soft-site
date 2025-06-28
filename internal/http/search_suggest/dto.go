package search_suggest

type SearchSuggestRequest struct {
	Query string `json:"query"`
}

type SearchSuggestResponsePayload struct {
	Items []SearchSuggestDTO `json:"items"`
}

type SearchSuggestDTO struct {
	SuggestType         string   `json:"suggest_type"`
	BannerTitle         *string  `json:"banner_title,omitempty"`
	BannerDescription   *string  `json:"banner_description,omitempty"`
	BannerImage         *string  `json:"banner_image,omitempty"`
	BannerURL           *string  `json:"banner_url,omitempty"`
	ItemID              *int64   `json:"item_id,omitempty"`
	ItemCategoryID      *int64   `json:"item_category_id,omitempty"`
	ItemTitle           *string  `json:"item_title,omitempty"`
	ItemCurrentPrice    *float64 `json:"item_current_price,omitempty"`
	ItemIsForSale       *bool    `json:"item_is_for_sale,omitempty"`
	ItemOldPrice        *float64 `json:"item_old_price,omitempty"`
	ItemThumbnailURL    *string  `json:"item_thumbnail_url,omitempty"`
	ItemType            *string  `json:"item_type,omitempty"`
	ItemHorizontalImage *string  `json:"item_horizontal_image,omitempty"`
	ItemGenres          []string `json:"item_genres,omitempty"`
	ItemReleaseDate     *string  `json:"item_release_date,omitempty"`
	ItemInSub           *bool    `json:"item_in_sub"`
	Probability         float64  `json:"probability"`
}
