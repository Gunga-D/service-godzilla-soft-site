package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	deepseek "github.com/cohesion-org/deepseek-go"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type TopicAIResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

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
	topicThemesPromt string = "Представь, что ты продюсер статей для сайта, на котором продаются цифровые товары: " +
		"ключи Steam/EA/Microsoft и другие, подарки Steam, случайные игры, услуги пополнения игровых аккаунтов. " +
		"Сгенерируй список возможных статей, которые могут быть представлены на этом сайте. " +
		"Напиши их название и коротко опиши их содержание в виде оглавления. " +
		"Не акцентируй внимание на продаже игр в нашем магазине, а скорее предлагай статьи, которые могут заинтеровать игроков чем-то новым."
	topicGenerationPromt string = "Исходя из предложенных тобой варинатов статей, выбери один вариант и создай описание этой статьи в формате Json, объем статьи должен быть в рамках 2000-5000 символов."
)

func GenerateTopic(client *deepseek.Client) (TopicAIResponse, error) {
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
	ctx := context.Background()
	response, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return TopicAIResponse{}, errors.New(fmt.Sprintf("Themes generation response error: %v", err))
	}

	// generate topic according to themes
	request = &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekReasoner,
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
		return TopicAIResponse{}, errors.New(fmt.Sprintf("Topic generation response error: %v", err))
	}

	var resp TopicAIResponse
	err = json.Unmarshal([]byte(response.Choices[len(response.Choices)-1].Message.Content), &resp)
	if err != nil {
		return TopicAIResponse{}, errors.New(fmt.Sprintf("Unmarshal error: %v", err))
	}
	return resp, nil
}

func main() {
	// Set up the Deepseek client
	var index atomic.Int32
	var mu sync.Mutex

	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.FromSlash(workdir + "/topics/generated/")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	client := deepseek.NewClient(GetApiKey(), GetApiURL())

	var wg sync.WaitGroup
	errChan := make(chan error, 10) // buffer for all potential errors

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			topic, err := GenerateTopic(client)
			if err != nil {
				errChan <- err
				return
			}

			// Get and increment index atomically
			currentIndex := index.Add(1) - 1

			mu.Lock()
			defer mu.Unlock()

			path := filepath.FromSlash(dir + fmt.Sprintf("topic%d.md", currentIndex))
			file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			_, err = file.Write([]byte(topic.Content))
			if err != nil {
				errChan <- err
			}
		}()
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Print any errors that occurred
	for err := range errChan {
		fmt.Println("Error:", err)
	}
}
