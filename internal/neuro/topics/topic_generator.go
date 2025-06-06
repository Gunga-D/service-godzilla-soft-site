package topics

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	topics "github.com/Gunga-D/service-godzilla-soft-site/internal/topics/postgres"
	"github.com/cohesion-org/deepseek-go"
	"log"
	"os"
	"time"
)

func GetApiKey() string {
	res := os.Getenv("DEEPSEEK_TOKEN")
	if len(res) <= 0 {
		log.Fatal("DEEPSEEK_TOKEN environment variable not set")
	}
	return res
}

func GetApiURL() string {
	res := os.Getenv("DEEPSEEK_URL")
	if len(res) <= 0 {
		log.Fatal("DEEPSEEK_URL environment variable not set")
	}
	return res
}

const (
	topicThemesPromt string = "Представь, что ты продюсер статей для сайта, на котором продаются различные цифровые товары и услуги: " +
		"ключи для игр Steam/EA/Microsoft и других платформ, рулетка случайных игр, услуги пополнения игровых аккаунтов, услуга генерации списка игр на основе по предпочтениям пользователя на основе нейронной сети." +
		"Сгенерируй список возможных статей, которые могут быть представлены на этом сайте. " +
		"Напиши их название и коротко опиши их содержание в виде оглавления. " +
		"Не акцентируй внимание на продаже игр в нашем магазине, а скорее предлагай статьи, которые могут заинтеровать игроков чем-то новым."
	topicGenerationPromt string = "Исходя из предложенных тобой варинатов статей, выбери один вариант и создай описание этой статьи в формате Json. " +
		"Объем статьи должен быть в рамках 2000-5000 символов." +
		"Главный приоритет сгенерированный статьи - часто выпадать в поисковых запросах, поэтому создавай статьи, которые оптимизированы под SEO"
)

func GenerateTopic(ctx context.Context, client *deepseek.Client) (topics.Topic, error) {
	// generate topic themes
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekReasoner,
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role: deepseek.ChatMessageRoleUser, Content: topicThemesPromt,
			},
		},
	}

	// Send the request and handle the response
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return topics.Topic{}, errors.New(fmt.Sprintf("Themes generation response error: %v", err))
	}

	// generate topic according to themes
	request = &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "json_object",
		},
		Messages: []deepseek.ChatCompletionMessage{
			request.Messages[0], // previous user message
			{
				Role: deepseek.ChatMessageRoleAssistant, Content: response.Choices[0].Message.Content, // assistant response
			},
			{
				Role: deepseek.ChatMessageRoleSystem, Content: "Отвечать на запросы нужно в строгом Json формате:" +
					"{ 'title' : 'Заголовок статьи' }" +
					"{ 'content' : 'Основная часть статьи в формате Markdown' }",
			},
			{
				Role: deepseek.ChatMessageRoleUser, Content: topicGenerationPromt,
			},
		},
	}

	// wait for specific topic response
	response, err = client.CreateChatCompletion(ctx, request)
	if err != nil {
		return topics.Topic{}, errors.New(fmt.Sprintf("Topic generation response error: %v", err))
	}

	// fill responce struct and return
	var resp topics.Topic
	err = json.Unmarshal([]byte(response.Choices[len(response.Choices)-1].Message.Content), &resp)
	if err != nil {
		return topics.Topic{}, errors.New(fmt.Sprintf("Unmarshal error: %v", err))
	}

	resp.CreatedAt = time.Now()
	return resp, nil
}
