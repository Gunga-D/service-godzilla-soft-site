package admin_save_thumbnail

type AdminSaveThumbnailRequest struct {
	FileName   string `json:"file_name"`
	DataBase64 string `json:"data_base64"`
}
