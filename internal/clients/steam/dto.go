package steam

type ResolveProfileIDResponse struct {
	Response struct {
		Success int     `json:"success"`
		Message *string `json:"message,omitempty"`
		SteamID *string `json:"steamid,omitempty"`
	} `json:"response"`
}

type GetProfileInfoResponse []ProfileInfo

type ProfileInfo struct {
	SteamID         string `json:"steamid"`
	AccountID       int64  `json:"accountid"`
	PersonaName     string `json:"persona_name"`
	AvatarUrl       string `json:"avatar_url"`
	ProfileUrl      string `json:"profile_url"`
	PersonaState    int    `json:"persona_state"`
	City            string `json:"city"`
	State           string `json:"state"`
	Country         string `json:"country"`
	RealName        string `json:"real_name,omitempty"`
	IsFriend        *bool  `json:"is_friend,omitempty"`
	FriendsInCommon *int   `json:"friends_in_common,omitempty"`
}

type AppDetailsResponse map[string]AppDetailsResult

type AppDetailsResult struct {
	Success bool       `json:"success"`
	Data    AppDetails `json:"data"`
}

type AppDetails struct {
	Type                string `json:"type"`
	Name                string `json:"name"`
	SteamAppid          int64  `json:"steam_appid"`
	IsFree              bool   `json:"is_free"`
	DetailedDescription string `json:"detailed_description"`
	AboutTheGame        string `json:"about_the_game"`
	ShortDescription    string `json:"short_description"`
	SupportedLanguages  string `json:"supported_languages"`
	HeaderImage         string `json:"header_image"`
	CapsuleImage        string `json:"capsule_image"`
	CapsuleImagev5      string `json:"capsule_imagev5"`
	PcRequirements      struct {
		Minimum     *string `json:"minimum,omitempty"`
		Recommended *string `json:"recommended,omitempty"`
	} `json:"pc_requirements"`
	MacRequirements struct {
		Minimum     *string `json:"minimum,omitempty"`
		Recommended *string `json:"recommended,omitempty"`
	} `json:"mac_requirements"`
	LinuxRequirements struct {
		Minimum     *string `json:"minimum,omitempty"`
		Recommended *string `json:"recommended,omitempty"`
	} `json:"linux_requirements"`
	Developers    []string `json:"developers"`
	Publishers    []string `json:"publishers"`
	PriceOverview struct {
		Currency         string  `json:"currency"`
		Initial          int64   `json:"initial"`
		Final            int64   `json:"final"`
		DiscountPercent  float64 `json:"discount_percent"`
		InitialFormatted string  `json:"initial_formatted"`
		FinalFormatted   string  `json:"final_formatted"`
	} `json:"price_overview"`
	Platforms struct {
		Windows bool `json:"windows"`
		Mac     bool `json:"mac"`
		Linux   bool `json:"linux"`
	} `json:"platforms"`
	Categories []struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"categories"`
	Genres []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"genres"`
	Screenshots []struct {
		ID            int    `json:"id"`
		PathThumbnail string `json:"path_thumbnail"`
		PathFull      string `json:"path_full"`
	} `json:"screenshots"`
	Movies []struct {
		ID        int    `json:"movies"`
		Name      string `json:"name"`
		Thumbnail string `json:"thumbnail"`
		Webm      struct {
			Res480 string `json:"480"`
			ResMax string `json:"max"`
		} `json:"webm"`
		MP4 struct {
			Res480 string `json:"480"`
			ResMax string `json:"max"`
		} `json:"mp4"`
		Highlight bool `json:"highlight"`
	} `json:"movies"`
	ReleaseDate struct {
		ComingSoon bool   `json:"coming_soon"`
		Date       string `json:"date"`
	} `json:"release_date"`
	Background string `json:"background"`
}
