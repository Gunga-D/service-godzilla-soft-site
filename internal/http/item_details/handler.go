package item_details

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/AlekSi/pointer"
	api "github.com/Gunga-D/service-godzilla-soft-site/internal/http"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/logger"
)

type handler struct {
	itemGetter    itemGetter
	recomendation recomendation
}

func NewHandler(itemGetter itemGetter, recomendation recomendation) *handler {
	return &handler{
		itemGetter:    itemGetter,
		recomendation: recomendation,
	}
}

func (h *handler) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Query().Get("item_id"), 10, 64)
		if err != nil {
			api.Return400("Невалидный запрос", w)
			return
		}
		item, err := h.itemGetter.GetItemByID(r.Context(), id)
		if err != nil {
			api.Return500("Неизвестная ошибка", w)
			return
		}
		if item == nil {
			api.Return404("Такого товара нет в наличии", w)
			return
		}

		wg := sync.WaitGroup{}

		var similarGames []SimilarGameDTO
		wg.Add(1)
		go func() {
			defer wg.Done()
			if item.SteamBlock == nil {
				return
			}
			recommendedItems, err := h.recomendation.Recommend(r.Context(), item.ID, item.SteamBlock.Genres)
			if err != nil {
				log.Printf("error to get recommended items: %v\n", err)
				return
			}
			for _, recItem := range recommendedItems {
				itemType := "cdkey"
				if recItem.IsSteamGift {
					itemType = "gift"
				}

				similarGames = append(similarGames, SimilarGameDTO{
					ID:           recItem.ID,
					Type:         itemType,
					Title:        recItem.Title,
					CategoryID:   recItem.CategoryID,
					ThumbnailURL: recItem.ThumbnailURL,
					CurrentPrice: float64(recItem.CurrentPrice) / 100,
				})
			}
		}()

		var oldPrice *float64
		if item.OldPrice != nil {
			oldPrice = pointer.ToFloat64(float64(*item.OldPrice) / 100)
		}

		itemType := "cdkey"
		if item.IsSteamGift {
			itemType = "gift"
		}

		var bxGalleryUrls []string
		if item.SteamBlock != nil {
			for _, v := range item.SteamBlock.Screenshots {
				bxGalleryUrls = append(bxGalleryUrls, v.PathThumbnail)
			}
		}

		var bxMovies []MovieDTO
		if item.SteamBlock != nil {
			for _, v := range item.SteamBlock.Movies {
				videoNorm, err := toHTTPS(v.MP4.Res480)
				if err != nil {
					continue
				}
				bxMovies = append(bxMovies, MovieDTO{
					Poster: v.Thumbnail,
					Video:  videoNorm,
				})
			}
		}

		var creator *string
		if item.SteamBlock != nil {
			if len(item.SteamBlock.Developers) > 0 {
				creator = pointer.ToString(strings.Join(item.SteamBlock.Developers, ", "))
			}
		}

		var publisher *string
		if item.SteamBlock != nil {
			if len(item.SteamBlock.Publishers) > 0 {
				publisher = pointer.ToString(strings.Join(item.SteamBlock.Publishers, ", "))
			}
		}

		var releaseDate *string
		if item.SteamBlock != nil {
			releaseDate = pointer.ToString(item.SteamBlock.ReleaseDate)
		}

		desc := item.Description
		if desc == nil && item.SteamBlock != nil {
			desc = pointer.ToString(item.SteamBlock.ShortDescription)
		}

		var pcRequirements *SteamRequirementsDTO
		if item.SteamBlock != nil {
			pcRequirements = &SteamRequirementsDTO{
				Minimum:     item.SteamBlock.PcRequirements.Minimum,
				Recommended: item.SteamBlock.PcRequirements.Recommended,
			}
		}

		var bxImageURL *string
		if item.SteamAppID != nil {
			bxImageURL = pointer.ToString(fmt.Sprintf("https://steamcdn-a.akamaihd.net/steam/apps/%d/library_600x900.jpg", *item.SteamAppID))
		}

		var genres []string
		if item.SteamBlock != nil {
			genres = item.SteamBlock.Genres
		}

		backgroundURl := item.BackgroundURL
		if backgroundURl == nil && item.SteamBlock != nil {
			backgroundURl = &item.SteamBlock.Background
		}

		var yandexMarketBlock *YandexMarketDTO
		if item.YandexMarket != nil {
			yandexMarketBlock = &YandexMarketDTO{
				Price:        item.YandexMarket.Price,
				Rating:       item.YandexMarket.Rating,
				ReviewsCount: item.YandexMarket.ReviewsCount,
			}
		}

		wg.Wait()
		itemDTO := ItemDTO{
			ID:             item.ID,
			Title:          item.Title,
			Type:           itemType,
			Description:    desc,
			CategoryID:     item.CategoryID,
			Platform:       item.Platform,
			Region:         item.Region,
			Publisher:      publisher,
			Creator:        creator,
			ReleaseDate:    releaseDate,
			CurrentPrice:   float64(item.CurrentPrice) / 100,
			IsForSale:      item.IsForSale,
			OldPrice:       oldPrice,
			ThumbnailURL:   item.ThumbnailURL,
			BackgroundURL:  backgroundURl,
			BxImageURL:     bxImageURL,
			BxGalleryUrls:  bxGalleryUrls,
			BxMovies:       bxMovies,
			Slip:           item.Slip,
			YandexMarket:   yandexMarketBlock,
			Genres:         genres,
			PcRequirements: pcRequirements,
			SimilarGames:   similarGames,
		}

		logger.Get().Log(fmt.Sprintf("👀 Товар\"%s\" просмотрели", item.Title))

		api.ReturnOK(itemDTO, w)
	}
}

func toHTTPS(addr string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	u.Scheme = "https"
	return u.String(), nil
}
