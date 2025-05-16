package steam

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

type filler struct{}

func NewFiller() *filler {
	return &filler{}
}

func (f *filler) Fill(ctx context.Context, items []item.ItemCache) error {
	for idx := 0; idx < len(items); idx++ {
		v := items[idx]
		if v.SteamRawData != nil && v.SteamAppID != nil {
			rawData, err := base64.StdEncoding.DecodeString(*v.SteamRawData)
			if err != nil {
				log.Printf("cannot decode steam app details of %d: %v", *v.SteamAppID, err)
				continue
			}
			var res steam.AppDetails
			if err := json.Unmarshal(rawData, &res); err != nil {
				log.Printf("cannot unmarshal steam app details of %d: %v", *v.SteamAppID, err)
				continue
			}

			screenshots := make([]item.SteamScreenshot, 0, len(res.Screenshots))
			for _, s := range res.Screenshots {
				screenshots = append(screenshots, item.SteamScreenshot{
					ID:            s.ID,
					PathThumbnail: s.PathThumbnail,
					PathFull:      s.PathFull,
				})
			}

			movies := make([]item.SteamMovie, 0, len(res.Movies))
			for _, m := range res.Movies {
				movies = append(movies, item.SteamMovie{
					ID:        m.ID,
					Name:      m.Name,
					Thumbnail: m.Thumbnail,
					Webm: item.SteamMovieFormat{
						Res480: m.Webm.Res480,
						ResMax: m.Webm.ResMax,
					},
					MP4: item.SteamMovieFormat{
						Res480: m.MP4.Res480,
						ResMax: m.MP4.ResMax,
					},
					Highlight: m.Highlight,
				})
			}

			genres := make([]string, 0, len(res.Genres))
			for _, g := range res.Genres {
				genres = append(genres, g.Description)
			}

			items[idx].SteamBlock = &item.ItemSteamBlock{
				DetailedDescription: res.DetailedDescription,
				AboutTheGame:        res.AboutTheGame,
				ShortDescription:    res.ShortDescription,
				HeaderImage:         res.HeaderImage,
				CapsuleImage:        res.CapsuleImage,
				CapsuleImagev5:      res.CapsuleImagev5,
				PcRequirements: item.SteamRequirements{
					Minimum:     res.PcRequirements.Minimum,
					Recommended: res.PcRequirements.Recommended,
				},
				Developers: res.Developers,
				Publishers: res.Publishers,
				Platforms: item.SteamPlatforms{
					Windows: res.Platforms.Windows,
					Mac:     res.Platforms.Mac,
					Linux:   res.Platforms.Linux,
				},
				Screenshots: screenshots,
				Movies:      movies,
				ReleaseDate: res.ReleaseDate.Date,
				Genres:      genres,
				Background:  res.Background,
			}
		}
	}
	return nil
}
