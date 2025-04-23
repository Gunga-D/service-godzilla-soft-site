package deepseek

import "context"

type Client interface {
	Completions(ctx context.Context, in CompletionsRequest) (*CompletionsResponse, error)
}
