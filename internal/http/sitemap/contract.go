package sitemap

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
)

type getter interface {
	FetchAllItems(ctx context.Context) ([]item.ItemCache, error)
}

type topicSource interface {
	FetchAllTopics(ctx context.Context) ([]topics.Topic, error)
}
