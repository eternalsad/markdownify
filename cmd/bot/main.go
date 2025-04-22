package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/md2"
	"github.com/gomarkdown/markdown/parser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Конфигурация бота
type Config struct {
	BotToken string
	Offset   int
}

// Структуры для работы с API Telegram
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ID int `json:"id"`
}

type SendMessageRequest struct {
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type ApiResponse struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

// GetUpdates получает обновления от API Telegram
func GetUpdates(config *Config) ([]Update, error) {
	// Формируем URL для API запроса
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d&timeout=60", config.BotToken, config.Offset)

	// Отправляем GET запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и проверяем ответ
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Ok {
		return nil, fmt.Errorf("API вернул ошибку: %s", string(body))
	}

	var updates []Update
	if err := json.Unmarshal(apiResp.Result, &updates); err != nil {
		return nil, err
	}

	return updates, nil
}

// SendMarkdownMessage отправляет сообщение в формате Markdown V2
func SendMarkdownMessage(botToken string, chatID int, text string) error {
	// Формируем URL для API запроса
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	// Создаем данные для запроса
	reqData := SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}

	// Преобразуем данные в JSON
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	// Отправляем POST запрос
	resp, err := http.Post(url, "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API вернул статус %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func Render(doc ast.Node, renderer *md2.Renderer) []byte {
	var buf bytes.Buffer
	renderer.RenderHeader(&buf, doc)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return renderer.RenderNode(&buf, node, entering)
	})
	renderer.RenderFooter(&buf, doc)
	return buf.Bytes()
}

// Список символов, которые нужно экранировать в Markdown V2
var specialChars = []string{
	"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!",
}

// Правильно экранирует специальные символы для Markdown V2
func escapeMarkdownV2(text string) string {
	// Создаем копию текста для модификации
	result := text

	// Экранируем все специальные символы
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	return result
}

var printAst = true

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	if printAst {
		fmt.Print("--- AST tree:\n")
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")
	}

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func main() {
	//input, err := os.ReadFile("./cmd/bot/input.txt")
	//if err != nil {
	//	log.Fatalf("fail reading input txt %s", err)
	//}
	//
	//extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	//p := parser.NewWithExtensions(extensions)
	//
	//doc := p.Parse(input)
	//renderer := md2.NewRenderer()
	//
	//output := Render(doc, renderer)
	//fmt.Println("huyhuy")
	//fmt.Println(string(output))

	//output := mdToHTML(input)

	//os.WriteFile("output.txt", output, os.ModePerm)

	output, err := os.ReadFile("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Получаем токен бота из переменных окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Необходимо установить переменную окружения TELEGRAM_BOT_TOKEN")
	}

	config := &Config{
		BotToken: botToken,
		Offset:   0,
	}

	fmt.Println("Бот запущен. Ожидание сообщений...")

	// Бесконечный цикл для получения и обработки обновлений
	for {
		updates, err := GetUpdates(config)
		if err != nil {
			log.Printf("Ошибка при получении обновлений: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Обрабатываем каждое обновление
		for _, update := range updates {
			// Обновляем offset (чтобы не получать одни и те же сообщения повторно)
			config.Offset = update.UpdateID + 1

			// Проверяем, есть ли текст сообщения
			if update.Message.Text != "" {
				chatID := update.Message.Chat.ID
				log.Printf("Получено сообщение от chat_id %d: %s", chatID, update.Message.Text)

				// Отправляем демонстрационное сообщение в Markdown V2
				demoMessage := string(output)

				err := SendMarkdownMessage(config.BotToken, chatID, demoMessage)
				if err != nil {
					log.Printf("Ошибка при отправке сообщения: %v", err)
				} else {
					log.Printf("Отправлено форматированное сообщение в чат %d", chatID)
				}
			}
		}

		// Небольшая пауза между запросами (чтобы не перегружать API)
		time.Sleep(40 * time.Second)
	}
}
