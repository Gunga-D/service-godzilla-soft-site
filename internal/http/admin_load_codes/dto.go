package admin_load_codes

type AdminLoadCodeRequest struct {
	ItemID int64    `json:"item_id"`
	Codes  []string `json:"codes"`
}
