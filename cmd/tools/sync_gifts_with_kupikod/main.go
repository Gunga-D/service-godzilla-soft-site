package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/kupikod"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/mtspay"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
	sq "github.com/Masterminds/squirrel"
)

var unavailableKupikodSettings = map[string]string{
	"unavailable":         "1",
	"amount_games_weekly": "true",
	"is_dlc":              "0",
	"with_totals":         "false",
	"per_page":            "150",
	"page":                "0",
}

var newsKupikodSettings = map[string]string{
	"new":               "1",
	"is_dlc":            "0",
	"sold_games_weekly": "true",
	"with_totals":       "true",
	"per_page":          "50",
	"page":              "0",
}

var popularKupikodSettings = map[string]string{
	"popular":             "1",
	"order":               "ASC",
	"is_dlc":              "0",
	"per_page":            "150",
	"page":                "1",
	"amount_games_weekly": "true",
}

func main() {
	postgres := postgres.Get(context.Background())

	kupikodClient := kupikod.NewClient()
	resp, err := kupikodClient.FetchGames(context.Background(), popularKupikodSettings)
	if err != nil {
		log.Fatalf("error to get games from kupikod: %v", err)
	}

	steamClient := steam.NewClient(os.Getenv("STEAM_KEY"))
	mtspayClient := mtspay.NewClient()

	rateResp, err := mtspayClient.GetRate(context.Background(), map[string]string{
		"amount": "30000",
	})
	if err != nil {
		log.Fatalf("error to get rate from mts pay: %v", err)
	}
	kztV, err := strconv.ParseFloat(rateResp.KztTopup, 64)
	if err != nil {
		log.Fatalf("cannot parse kzt rate: %v", err)
	}
	rateKztToRub := float64(30_050) / kztV

	var selfItems []item.Item
	for _, kupikodItem := range resp.Data {
		steamAppID, err := strconv.ParseInt(kupikodItem.ExternalID, 10, 64)
		if err != nil {
			fmt.Printf("%s: error to get steam app id: %v", kupikodItem.Name, err)
			continue
		}

		appDetails, err := steamClient.AppDetails(context.Background(), steamAppID)
		if err != nil {
			fmt.Printf("%s: error to get app details: %v", kupikodItem.Name, err)
			continue
		}

		selfPriceRub := float64(appDetails.PriceOverview.Final) * rateKztToRub / 100
		priceRub := int64((selfPriceRub * 1.134) * 100)
		limitRub := int64((selfPriceRub * 1.064) * 100)

		selfItems = append(selfItems, item.Item{
			Title:        kupikodItem.Name,
			SteamAppID:   &steamAppID,
			ThumbnailURL: appDetails.HeaderImage,
			CategoryID:   10001,
			Platform:     "Steam",
			Region:       "РФ",
			CurrentPrice: priceRub,
			IsForSale:    false,
			LimitPrice:   &limitRub,
			Status:       "active",
			Slip:         "1) После оплаты к вам в друзья Steam добавится бот. Пожалуйста, примите его в друзья, чтобы получить игру.\n2) После добавления, вам будет отправлен подарок в виде игры - нужно принять его.\n3) Игра в вашей библиотеке, можно играть",
			IsSteamGift:  true,
		})
	}

	// Добавление товара в базу
	q := sq.Insert("public.item").
		Columns(
			"title",
			"description",
			"category_id",
			"platform",
			"region",
			"current_price",
			"is_for_sale",
			"old_price",
			"thumbnail_url",
			"background_url",
			"status",
			"slip",
			"created_at",
			"updated_at",
			"steam_app_id",
			"is_steam_gift",
		)
	for idx := 0; idx < len(selfItems); idx++ {
		i := selfItems[idx]
		log.Printf("Добавили игру \"%s\"\n", i.Title)
		q = q.Values(
			i.Title,
			i.Description,
			i.CategoryID,
			i.Platform,
			i.Region,
			i.CurrentPrice,
			i.IsForSale,
			i.OldPrice,
			i.ThumbnailURL,
			i.BackgroundURL,
			i.Status,
			i.Slip,
			time.Now(),
			time.Now(),
			i.SteamAppID,
			i.IsSteamGift,
		)
	}
	query, args, err := q.
		Suffix(`
			ON CONFLICT (title) DO NOTHING;
		`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Fatalf("cannot create quary of inserting items: %v", err)
	}

	if _, err := postgres.Do(context.Background()).ExecContext(context.Background(), query, args...); err != nil {
		log.Fatalf("cannot insert all items: %v", err)
	}
}
