package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
)

func main() {
	c := steam.NewClient(os.Getenv("STEAM_KEY"))
	foundedApps, err := c.Search(context.Background(), "DOOM: The Dark Ages")
	if err != nil {
		log.Fatal(err)
	}
	appID, err := strconv.ParseInt((*foundedApps)[0].AppID, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found item - %d\n", appID)
	appDetails, err := c.AppDetails(context.Background(), appID)
	if err != nil {
		log.Fatal(err)
	}

	rawSteamData, err := json.Marshal(*appDetails)
	if err != nil {
		log.Fatal(err)
	}

	steamData := base64.StdEncoding.EncodeToString(rawSteamData)
	fmt.Println(steamData)
}
