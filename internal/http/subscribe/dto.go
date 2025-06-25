package subscribe

type SubscribeRequest struct {
	Email  *string `json:"email"`
	Period string  `json:"period"` // month, year
}

type SubscribeResponse struct {
	UserAccessToken *string `json:"user_access_token"`
	SubscriptionID  string  `json:"subscription_id"`
	PaymentLink     string  `json:"payment_link"`
}
