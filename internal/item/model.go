package item

import (
	"time"
)

const (
	ActiveStatus   = "active"
	PausedStatus   = "paused"
	ArchivedStatus = "archived"
)

type Item struct {
	ID          int64   `db:"id"`
	Title       string  `db:"title"`
	Description *string `db:"description"`
	CategoryID  int64   `db:"category_id"`
	Platform    string  `db:"platform"`
	Region      string  `db:"region"`
	// Указывается с копейками, таким образом:
	// 100 рублей = 10000
	CurrentPrice    int64     `db:"current_price"`
	IsForSale       bool      `db:"is_for_sale"`
	OldPrice        *int64    `db:"old_price"`
	LimitPrice      *int64    `db:"limit_price"`
	ThumbnailURL    string    `db:"thumbnail_url"`
	BackgroundURL   *string   `db:"background_url"`
	Status          string    `db:"status"`
	Slip            string    `db:"slip"`
	IsSteamGift     bool      `db:"is_steam_gift"`
	YandexID        *string   `db:"yandex_id"`
	SteamAppID      *int64    `db:"steam_app_id"`
	PriceLoc        *string   `db:"price_loc"`
	Popular         *int      `db:"popular"`
	New             *int      `db:"new"`
	Unavailable     bool      `db:"unavailable"`
	HorizontalImage *string   `db:"horizontal_image"`
	VerticalImage   *string   `db:"vertical_image"`
	SteamRawData    *string   `db:"steam_raw_data"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type ItemYandexMarketBlock struct {
	Price        float64
	ReviewsCount int
	Rating       float64
}

type ItemSteamBlock struct {
	DetailedDescription string
	AboutTheGame        string
	ShortDescription    string
	HeaderImage         string
	CapsuleImage        string
	CapsuleImagev5      string
	PcRequirements      SteamRequirements
	Developers          []string
	Publishers          []string
	Platforms           SteamPlatforms
	Screenshots         []SteamScreenshot
	Movies              []SteamMovie
	Genres              []string
	ReleaseDate         string
	Background          string
}

type SteamPlatforms struct {
	Windows bool
	Mac     bool
	Linux   bool
}

type SteamRequirements struct {
	Minimum     *string
	Recommended *string
}

type SteamScreenshot struct {
	ID            int
	PathThumbnail string
	PathFull      string
}

type SteamMovie struct {
	ID        int
	Name      string
	Thumbnail string
	Webm      SteamMovieFormat
	MP4       SteamMovieFormat
	Highlight bool
}

type SteamMovieFormat struct {
	Res480 string
	ResMax string
}

type ItemCache struct {
	Item
	YandexMarket *ItemYandexMarketBlock
	SteamBlock   *ItemSteamBlock
}
