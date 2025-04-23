package deepseek

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type client struct {
	rc *resty.Client
}

func NewClient(url string, token string) *client {
	if token == "" {
		log.Fatalf("deepseek: credentials must be non-empty")
	}

	rc := resty.New()
	rc.SetBaseURL(url)
	rc.Header.Set("User-Agent", "GodzillaSoft")
	rc.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return &client{
		rc: rc,
	}
}

func (c *client) Completions(ctx context.Context, in CompletionsRequest) (*CompletionsResponse, error) {
	resp, err := c.rc.SetTimeout(time.Minute).R().
		SetContext(ctx).
		SetBody(in).
		SetResult(CompletionsResponse{}).
		Post("/chat/completions")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	return resp.Result().(*CompletionsResponse), nil
}
