package yandex_market

import "context"

type Client interface {
	OfferMappings(ctx context.Context, req OfferMappingsRequest) (*OfferMappingsResponse, error)
	GoodsFeedback(ctx context.Context, req GoodsFeedbackRequest) (*GoodsFeedbackResponse, error)
}
