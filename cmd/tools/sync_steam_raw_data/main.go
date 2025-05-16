package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	item_postgres "github.com/Gunga-D/service-godzilla-soft-site/internal/item/postgres"
	"github.com/Gunga-D/service-godzilla-soft-site/pkg/postgres"
)

func main() {
	postgres := postgres.Get(context.Background())
	itemsRepo := item_postgres.NewRepo(postgres)
	items, err := itemsRepo.FetchItemsPaginatedCursorItemId(context.Background(), 1000, 0)
	if err != nil {
		log.Fatalf("cannot fetch items: %v", err)
	}

	steamClient := steam.NewClient(os.Getenv("STEAM_KEY"))
	for _, i := range items {
		if i.SteamAppID == nil {
			continue
		}

		appDetails, err := steamClient.AppDetails(context.Background(), *i.SteamAppID)
		if err != nil {
			log.Printf("cannot get steam app details - %d: %v\n", i.ID, err)
			continue
		}
		rawData, err := json.Marshal(*appDetails)
		if err != nil {
			log.Printf("cannot marshal steam app details - %d: %v\n", i.ID, err)
			continue
		}
		steamData := base64.StdEncoding.EncodeToString(rawData)

		err = itemsRepo.UpdateSteamRawData(context.Background(), i.ID, steamData)
		if err != nil {
			log.Printf("cannot update steam raw data - %d: %v\n", i.ID, err)
			continue
		}
		log.Printf("update steam raw data for item - %d\n", i.ID)
	}
}
