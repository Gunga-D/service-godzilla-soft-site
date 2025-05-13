package user_profile

type UserProfileResponsePayload struct {
	UserID    int64   `json:"user_id"`
	Email     *string `json:"email,omitempty"`
	SteamLink *string `json:"steam_link,omitempty"`
	PhotoURL  *string `json:"photo_url,omitempty"`
	Username  *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
}
