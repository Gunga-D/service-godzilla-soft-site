package steam

import "context"

type Client interface {
	ResolveProfileID(ctx context.Context, profileID string) (int64, error)
	GetProfileInfo(ctx context.Context, profileID int64) (*ProfileInfo, error)
	AppDetails(ctx context.Context, appID int64) (*AppDetails, error)
	GetGenreApps(ctx context.Context, genre string) (*GenreList, error)
	FetchPrices(ctx context.Context, appIds []string, loc *string) (*FetchPricesResponse, error)
}
