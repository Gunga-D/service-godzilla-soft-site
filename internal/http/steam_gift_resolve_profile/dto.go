package steam_gift_resolve_profile

type SteamGiftResolveProfileRequest struct {
	ProfileURL string `json:"profile_url"`
}

type SteamGiftResolveProfileResponse struct {
	AvatarURL   *string `json:"avatar_url,omitempty"`
	ProfileName string  `json:"profile_name"`
}
