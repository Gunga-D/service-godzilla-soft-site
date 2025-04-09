package fetch_collections

type CollectionDTO struct {
	ID              int64  `json:"id"`
	CategoryID      int64  `json:"category_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	BackgroundImage string `json:"background_image"`
}
