package collection_details

import (
	"context"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/collection"
)

type getter interface {
	GetCollectionByID(ctx context.Context, id int64) (*collection.Collection, error)
}
