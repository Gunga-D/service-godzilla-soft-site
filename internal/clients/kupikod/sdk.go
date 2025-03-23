package kupikod

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

func (c *client) FetchGames(ctx context.Context, queryParams map[string]string) (*FetchGamesResponse, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(queryParams).
		SetResult(FetchGamesResponse{}).
		Get("https://steam.kupikod.com/backend/api/games")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	return resp.Result().(*FetchGamesResponse), nil
}
