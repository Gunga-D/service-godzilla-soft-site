package steam

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

type client struct {
	rc    *resty.Client
	token string
}

func NewClient(token string) *client {
	if token == "" {
		log.Fatalf("steam invoice: credentials must be non-empty")
	}

	rc := resty.New()
	rc.Header.Set("User-Agent", "GodzillaSoft")

	return &client{
		rc:    rc,
		token: token,
	}
}

func (c *client) ResolveProfileID(ctx context.Context, profileID string) (int64, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"key":       c.token,
			"vanityurl": profileID,
		}).
		SetResult(ResolveProfileIDResponse{}).
		Get("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v1")
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	res := resp.Result().(*ResolveProfileIDResponse)
	if res.Response.Success != 1 {
		if res.Response.Message != nil {
			return 0, fmt.Errorf(*res.Response.Message)
		}
		return 0, fmt.Errorf("invalid success of response = %d", res.Response.Success)
	}
	if res.Response.SteamID == nil {
		return 0, fmt.Errorf("cannot resolve steam profile id")
	}
	steamID, err := strconv.ParseInt(*res.Response.SteamID, 10, 64)
	if err != nil {
		return 0, err
	}
	return steamID, nil
}

func (c *client) GetProfileInfo(ctx context.Context, profileID int64) (*ProfileInfo, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"steamids": fmt.Sprint(profileID),
		}).
		SetResult(GetProfileInfoResponse{}).
		Get("https://steamcommunity.com/actions/ajaxresolveusers")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	res := *resp.Result().(*GetProfileInfoResponse)
	return &res[0], nil
}

func (c *client) AppDetails(ctx context.Context, appID int64) (*AppDetails, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"appids": fmt.Sprint(appID),
			"cc":     "KZ",
			"l":      "russian",
		}).
		SetResult(AppDetailsResponse{}).
		Get("https://store.steampowered.com/api/appdetails")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	res := *resp.Result().(*AppDetailsResponse)
	appDetails := res[fmt.Sprint(appID)]
	if !appDetails.Success {
		return nil, fmt.Errorf("cannot get app")
	}
	return &appDetails.Data, nil
}

func (c *client) GetGenreApps(ctx context.Context, genre string) (*GenreList, error) {
	resp, err := c.rc.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"genre": genre,
			"cc":    "KZ",
			"l":     "russian",
		}).
		SetResult(GenreList{}).
		Get("https://store.steampowered.com/api/getappsingenre")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code is not ok = %d", resp.StatusCode())
	}
	return resp.Result().(*GenreList), nil
}
