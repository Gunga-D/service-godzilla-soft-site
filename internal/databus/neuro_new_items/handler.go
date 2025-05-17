package neuro_new_items

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/category"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/deepseek"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/clients/steam"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/databus"
	"github.com/Gunga-D/service-godzilla-soft-site/internal/item"
)

const (
	prompt       = "Подбери игры, которые необходимо купить на платформе Steam и выдай результат в формате (строго соблюдай данный формат): Название игр с разделительным символом \n.Не используй символы форматирования."
	rateKztToRub = 5.7
)

type handler struct {
	itemsRepo      item.Repository
	deepseekClient deepseek.Client
	steamClient    steam.Client
	consumer       databus.Consumer
}

func NewHandler(itemsRepo item.Repository, deepseekClient deepseek.Client, steamClient steam.Client, consumer databus.Consumer) *handler {
	return &handler{
		itemsRepo:      itemsRepo,
		deepseekClient: deepseekClient,
		steamClient:    steamClient,
		consumer:       consumer,
	}
}

func (h *handler) Consume(ctx context.Context) {
	msgs, err := h.consumer.ConsumeDatabusNeuroNewItems(ctx)
	if err != nil {
		log.Fatalf("cannot start consume databus neuro task: %v", err)
	}
	for msg := range msgs {
		var data databus.NeuroNewItemsDTO
		if err = json.Unmarshal(msg.Body, &data); err != nil {
			msg.Ack(false)
			continue
		}
		log.Printf("neuro new items: new query %s\n", data.Query)

		resp, err := h.deepseekClient.Completions(ctx, deepseek.CompletionsRequest{
			Model: "deepseek-chat",
			Messages: []deepseek.MessageDTO{
				{
					Role:    "system",
					Content: prompt,
				},
				{
					Role:    "user",
					Content: data.Query,
				},
			},
			Stream: false,
		})
		if err != nil {
			log.Printf("cannot neuro new items: deepseek completions err: %v\n", err)
			msg.Nack(false, true)
			continue
		}
		if len(resp.Choices) == 0 {
			log.Println("cannot neuro new items: no deepseek result")
			msg.Ack(false)
			continue
		}

		appNames := strings.Split(resp.Choices[0].Message.Content, "\n")
		for _, appName := range appNames {
			appResp, err := h.steamClient.Search(ctx, appName)
			if err != nil {
				log.Printf("neuro new items: skip steam appName - %s cause cannot find it: %v\n", appName, err)
				continue
			}
			if appResp == nil || len(*appResp) == 0 {
				log.Printf("neuro new items: skip steam appName - %s cause no response\n", appName)
				continue
			}
			steamAppID, err := strconv.ParseInt((*appResp)[0].AppID, 10, 64)
			if err != nil {
				log.Printf("neuro new items: skip steam appName - %s cause appID invalid\n", appName)
				continue
			}
			i, err := h.itemsRepo.GetItemBySteamAppID(ctx, steamAppID)
			if err != nil {
				log.Printf("neuro new items: skip steam appName - %s cause cannot get items repo info: %v\n", appName, err)
				continue
			}
			if i != nil {
				log.Printf("neuro new items: skip steam appName - %s cause item already exists\n", appName)
				continue
			}

			appDetails, err := h.steamClient.AppDetails(ctx, steamAppID)
			if err != nil {
				log.Printf("neuro new items: skip steam appName - %s cause cannot get app details: %v\n", appName, err)
				continue
			}
			rawSteamData, err := json.Marshal(*appDetails)
			if err != nil {
				log.Printf("neuro new items: skip steam appName - %s cause cannot marshal steam data: %v\n", appName, err)
				continue
			}

			if appDetails.PriceOverview.Final == 0 {
				log.Printf("neuro new items: skip steam appName - %s cause game is free:\n", appName)
				continue
			}
			kzInRubPrice := float64(appDetails.PriceOverview.Final) / rateKztToRub
			selfPrice := kzInRubPrice / 100

			steamData := base64.StdEncoding.EncodeToString(rawSteamData)
			_, err = h.itemsRepo.CreateItem(ctx, item.Item{
				Title:        appDetails.Name,
				CategoryID:   category.GamesCategoryID,
				Platform:     "Steam",
				Region:       "РФ и СНГ",
				CurrentPrice: int64(selfPrice*1.175) * 100,
				IsForSale:    false,
				LimitPrice:   pointer.ToInt64(int64(selfPrice*1.084) * 100),
				ThumbnailURL: "https://disk.godzillasoft.ru/no_image_placeholder_v1.png",
				Status:       "active",
				Slip: `<ol class='BxItemInstruction'>
	<li>После оплаты к вам в друзья Steam добавится бот. Пожалуйста, примите его в друзья, чтобы получить игру.</li>
	<li>После добавления, вам будет отправлен подарок в виде игры - нужно принять его.</li>
	<li>Игра в вашей библиотеке, можно играть</li>
</ol>`,
				IsSteamGift:  true,
				SteamAppID:   &steamAppID,
				PriceLoc:     pointer.ToString("KZ"),
				SteamRawData: &steamData,
			})
			if err != nil {
				log.Printf("neuro new items: skip steam appName - %s cause cannot create new item: %v\n", appName, err)
				continue
			}
			log.Printf("neuro new items: added %s\n", appDetails.Name)
		}
		msg.Ack(false)
	}
}
