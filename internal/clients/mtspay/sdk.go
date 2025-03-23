package mtspay

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type client struct {
	rc *resty.Client
}

func NewClient() *client {
	rc := resty.New()
	return &client{
		rc: rc,
	}
}

func (c *client) GetRate(ctx context.Context, queryParams map[string]string) (*GetRateResponse, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParams).
		SetResult(GetRateResponse{}).
		Get("https://keys.foreignpay.ru/webhook/v2/topup/check_rate")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	return resp.Result().(*GetRateResponse), nil
}
