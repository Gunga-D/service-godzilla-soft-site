package kupikod

type FetchGamesResponse struct {
	Data []KupikodItemDTO `json:"data"`
}

type KupikodItemDTO struct {
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
}
