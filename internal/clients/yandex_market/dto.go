package yandex_market

type OfferMappingsRequest struct {
	OfferIds []string `json:"offerIds"`
}

type OfferMappingsResponse struct {
	Status string `json:"status"`
	Result struct {
		OfferMappings []struct {
			Offer struct {
				OfferId    string   `json:"offerId"`
				Name       string   `json:"name"`
				Pictures   []string `json:"pictures"`
				BasicPrice struct {
					Value        float64 `json:"value"`
					CurrencyId   string  `json:"currencyId"`
					DiscountBase float64 `json:"discountBase"`
				} `json:"basicPrice"`
			} `json:"offer"`
			Mapping struct {
				MarketModelId int64 `json:"marketModelId"`
			} `json:"mapping"`
		} `json:"offerMappings"`
	} `json:"result"`
}

type GoodsFeedbackRequest struct {
	ModelIds []int64 `json:"modelIds"`
}

type GoodsFeedbackResponse struct {
	Status string `json:"status"`
	Result struct {
		Feedbacks []struct {
			FeedbackId  int64  `json:"feedbackId"`
			CreatedAt   string `json:"createdAt"`
			Author      string `json:"author"`
			Description struct {
				Advantages    string `json:"advantages"`
				Disadvantages string `json:"disadvantages"`
				Comment       string `json:"comment"`
			} `json:"description"`
			Identifiers struct {
				ModelID int64 `json:"modelId"`
			} `json:"identifiers"`
			Statistics struct {
				Rating int `json:"rating"`
			} `json:"statistics"`
		} `json:"feedbacks"`
		Paging struct {
			NextPageToken *string `json:"nextPageToken,omitempty"`
		} `json:"paging"`
	} `json:"result"`
}
