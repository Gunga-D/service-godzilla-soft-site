package item_details

type ItemDTO struct {
	ID                  int64                 `json:"id"`
	Type                string                `json:"type"`
	Title               string                `json:"title"`
	Description         *string               `json:"description"`
	DetailedDescription *string               `json:"detailed_description"`
	CategoryID          int64                 `json:"category_id"`
	Platform            string                `json:"platform"`
	Region              string                `json:"region"`
	Publisher           *string               `json:"publisher,omitempty"`
	Creator             *string               `json:"creator,omitempty"`
	ReleaseDate         *string               `json:"release_date,omitempty"`
	CurrentPrice        float64               `json:"current_price"`
	IsForSale           bool                  `json:"is_for_sale"`
	OldPrice            *float64              `json:"old_price"`
	ThumbnailURL        string                `json:"thumbnail_url"`
	BackgroundURL       *string               `json:"background_url,omitempty"`
	BxImageURL          *string               `json:"bx_image_url,omitempty"`
	BxGalleryUrls       []string              `json:"bx_gallery_urls,omitempty"`
	BxMovies            []MovieDTO            `json:"movies,omitempty"`
	PcRequirements      *SteamRequirementsDTO `json:"pc_requirements,omitempty"`
	Genres              []string              `json:"genres,omitempty"`
	Slip                string                `json:"slip"`
	YandexMarket        *YandexMarketDTO      `json:"yandex_market,omitempty"`
	SimilarGames        []SimilarGameDTO      `json:"similar_games,omitempty"`
}

type MovieDTO struct {
	Poster string `json:"poster"`
	Video  string `json:"video"`
}

type SteamRequirementsDTO struct {
	Minimum     *string `json:"minimun"`
	Recommended *string `json:"recommended"`
}

type YandexMarketDTO struct {
	Rating       float64 `json:"rating"`
	Price        float64 `json:"price"`
	ReviewsCount int     `json:"reviews_count"`
}

type SimilarGameDTO struct {
	ID           int64   `json:"id"`
	Type         string  `json:"type"`
	Title        string  `json:"title"`
	CategoryID   int64   `json:"category_id"`
	ThumbnailURL string  `json:"thumbnail_url"`
	CurrentPrice float64 `json:"current_price"`
}
