package yandex_market

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type client struct {
	rc         *resty.Client
	businessID int64
}

func NewClient(apiURL string, authHeader string, businessID int64) *client {
	if apiURL == "" || authHeader == "" || businessID == 0 {
		log.Fatalf("yandexMarket: credentials must be non-empty")
	}

	rc := resty.New()
	rc.SetBaseURL(apiURL)
	rc.Header.Set("User-Agent", "GodzillaSoft")
	rc.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authHeader))

	return &client{
		rc:         rc,
		businessID: businessID,
	}
}

func (c *client) OfferMappings(ctx context.Context, req OfferMappingsRequest) (*OfferMappingsResponse, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(OfferMappingsResponse{}).
		Post(fmt.Sprintf("/businesses/%d/offer-mappings", c.businessID))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("cannot get offer mappings: \n%s", string(resp.Body()))
	}
	return resp.Result().(*OfferMappingsResponse), nil
}

func (c *client) GoodsFeedback(ctx context.Context, req GoodsFeedbackRequest) (*GoodsFeedbackResponse, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(GoodsFeedbackResponse{}).
		Post(fmt.Sprintf("/businesses/%d/goods-feedback?limit=300", c.businessID))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("cannot get goods feedback: \n%s", string(resp.Body()))
	}
	return resp.Result().(*GoodsFeedbackResponse), nil
}
