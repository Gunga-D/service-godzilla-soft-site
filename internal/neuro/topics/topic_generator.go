package topics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	topics "github.com/Gunga-D/service-godzilla-soft-site/internal/topics"
	"github.com/cohesion-org/deepseek-go"
)

func GetApiKey() string {
	res := os.Getenv("DEEPSEEK_TOKEN")
	if len(res) <= 0 {
		log.Fatal("DEEPSEEK_TOKEN environment variable not set")
	}
	return res
}

func GetApiURL() string {
	res := os.Getenv("DEEPSEEK_TOPIC_URL")
	if len(res) <= 0 {
		log.Fatal("DEEPSEEK_URL environment variable not set")
	}
	return res
}

const (
	topicThemesPromt string = "Представь, что ты продюсер статей для сайта, на котором продаются различные цифровые товары и услуги: " +
		"ключи для игр Steam/EA/Microsoft и других платформ, рулетка случайных игр, услуги пополнения игровых аккаунтов, услуга генерации списка игр на основе по предпочтениям пользователя на основе нейронной сети." +
		"Сгенерируй список возможных статей, которые могут быть представлены на этом сайте." +
		"Напиши их название и коротко опиши их содержание в виде оглавления." +
		"Не акцентируй внимание на продаже игр в нашем магазине, а скорее предлагай статьи, которые могут заинтеровать игроков чем-то новым."
	topicGenerationPromt = "Представь, что ты продюсер статей для сайта, на котором продаются различные цифровые товары и услуги: " +
		"ключи для игр Steam/EA/Microsoft и других платформ, рулетка случайных игр, услуги пополнения игровых аккаунтов, услуга генерации списка игр на основе по предпочтениям пользователя на основе нейронной сети." +
		"Твоя задача написать статью по выбранной теме как пример публикации на таком сайте." +
		"Объем статьи должен быть в рамках 2000-5000 символов." +
		"Главный приоритет сгенерированный статьи - часто выпадать в поисковых запросах, поэтому создавай статьи, которые оптимизированы под SEO" +
		"Тема статьи должна быть следующая: %s." +
		"Тело статьи должно содержать выдерживать следующую структуру: %s."
)

type Theme struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Response struct {
	Themes []Theme `json:"articles"`
}

func makeTopicPromt(theme Theme) string {
	return fmt.Sprintf(topicGenerationPromt, theme.Title, theme.Content)
}

func GenerateThemes(ctx context.Context, client *deepseek.Client) (Response, error) {
	// generate topic themes
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekReasoner,
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "json_object",
		},
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role: deepseek.ChatMessageRoleUser, Content: topicThemesPromt,
			},
			{
				Role: deepseek.ChatMessageRoleSystem, Content: "Ответ должен содержать только сплошное описание в Json-формате." +
					"Объект `articles` должен содержать в себе элементы с ключом title и содержание с ключом content.",
			},
		},
	}
	// Send the request and handle the response
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return Response{}, err
	}
	var resp Response
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &resp)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}

func GenerateTopic(ctx context.Context, client *deepseek.Client, theme Theme) (topics.Topic, error) {
	request := &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "json_object",
		},
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role: deepseek.ChatMessageRoleSystem, Content: "Отвечать на запросы нужно в строгом Json формате:" +
					"{ 'title' : 'Заголовок статьи' }" +
					"{ 'topic_content' : 'Основная часть статьи в формате Markdown' }",
			},
			{
				Role: deepseek.ChatMessageRoleUser, Content: makeTopicPromt(theme),
			},
		},
	}

	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return topics.Topic{}, err
	}

	var topic topics.Topic
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &topic)
	if err != nil {
		return topics.Topic{}, err
	}

	return topic, nil
}
