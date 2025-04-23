package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/md2"
	"github.com/gomarkdown/markdown/parser"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

// Render выполняет рендеринг AST документа в формат Markdown V2
func Render(doc ast.Node, renderer *md2.Renderer) []byte {
	var buf bytes.Buffer
	renderer.RenderHeader(&buf, doc)
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		return renderer.RenderNode(&buf, node, entering)
	})
	renderer.RenderFooter(&buf, doc)
	return buf.Bytes()
}

// Перерабатывает Markdown в формат Markdown V2 для Telegram
func processMarkdownFile(filePath string, printAst bool) (string, error) {
	// Читаем файл
	input, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла %s: %v", filePath, err)
	}

	// Создаем парсер Markdown с расширениями
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	// Парсим документ
	doc := p.Parse(input)

	// Если нужно, выводим AST дерево
	if printAst {
		fmt.Printf("--- AST tree для файла %s:\n", filePath)
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")
	}

	// Создаем рендерер для Markdown V2
	renderer := md2.NewRenderer()

	// Рендерим документ
	output := Render(doc, renderer)

	err = os.WriteFile("output.txt", output, os.ModePerm)
	if err != nil {
		log.Printf("output err: %s", err)
	}

	// Возвращаем результат как строку
	return string(output), nil
}

// Получает список тестовых файлов
func getTestFiles(testDirPath string) ([]string, error) {
	// Проверяем существование директории
	_, err := os.Stat(testDirPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("директория %s не существует", testDirPath)
	}

	// Получаем список файлов
	files, err := os.ReadDir(testDirPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения директории %s: %v", testDirPath, err)
	}

	// Собираем пути к файлам
	var filePaths []string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".md") || strings.HasSuffix(file.Name(), ".txt")) {
			filePaths = append(filePaths, filepath.Join(testDirPath, file.Name()))
		}
	}

	return filePaths, nil
}

// Отправляет все тестовые файлы пользователю
func sendAllTestFiles(botToken string, chatID int, testDirPath string, printAst bool) {
	// Получаем список файлов
	filePaths, err := getTestFiles(testDirPath)
	if err != nil {
		log.Printf("Ошибка при получении списка файлов: %v", err)
		SendMarkdownMessage(botToken, chatID, escapeMarkdownV2(fmt.Sprintf("Ошибка: %v", err)))
		return
	}

	// Если файлов нет, отправляем сообщение об этом
	if len(filePaths) == 0 {
		SendMarkdownMessage(botToken, chatID, escapeMarkdownV2("Тестовые файлы не найдены"))
		return
	}

	// Отправляем каждый файл
	for _, filePath := range filePaths {
		filename := filepath.Base(filePath)

		// Обрабатываем файл Markdown
		outputMarkdown, err := processMarkdownFile(filePath, printAst)
		if err != nil {
			log.Printf("Ошибка при обработке файла %s: %v", filename, err)
			SendMarkdownMessage(botToken, chatID, escapeMarkdownV2(fmt.Sprintf("Ошибка при обработке файла %s: %v", filename, err)))
			continue
		}

		// Отправляем сообщение с заголовком файла
		SendMarkdownMessage(botToken, chatID, escapeMarkdownV2(fmt.Sprintf("📁 Файл: %s", filename)))

		// Ждем немного, чтобы сообщения приходили в правильном порядке
		time.Sleep(500 * time.Millisecond)

		// Отправляем содержимое файла
		err = SendMarkdownMessage(botToken, chatID, outputMarkdown)
		if err != nil {
			log.Printf("Ошибка при отправке файла %s: %v", filename, err)
			SendMarkdownMessage(botToken, chatID, escapeMarkdownV2(fmt.Sprintf("Не удалось отправить содержимое файла %s: %v", filename, err)))
		} else {
			log.Printf("Файл %s успешно отправлен в чат %d", filename, chatID)
		}

		// Пауза между отправкой файлов, чтобы не превысить лимиты API
		time.Sleep(1 * time.Second)
	}

	// Отправляем сообщение о завершении
	SendMarkdownMessage(botToken, chatID, escapeMarkdownV2("✅ Все файлы отправлены"))
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

func main() {
	// Получаем токен бота из переменных окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Необходимо установить переменную окружения TELEGRAM_BOT_TOKEN")
	}

	// Путь к директории с тестовыми файлами
	testDirPath := "./cmd/bot/tests"

	// Создаем конфигурацию бота
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
				messageText := update.Message.Text
				log.Printf("Получено сообщение от chat_id %d: %s", chatID, messageText)

				// Если пользователь отправил команду /start или /files
				if messageText == "/start" || messageText == "/files" {
					// Отправляем приветственное сообщение
					SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2("Привет! Отправляю тестовые Markdown файлы..."))

					// Отправляем все тестовые файлы
					sendAllTestFiles(config.BotToken, chatID, testDirPath, true)
				} else if strings.HasPrefix(messageText, "/file ") {
					// Если пользователь запросил конкретный файл
					fileName := strings.TrimPrefix(messageText, "/file ")
					filePath := filepath.Join(testDirPath, fileName)

					// Проверяем существование файла
					if _, err := os.Stat(filePath); os.IsNotExist(err) {
						SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2(fmt.Sprintf("Файл %s не найден", fileName)))
						continue
					}

					// Обрабатываем и отправляем конкретный файл
					outputMarkdown, err := processMarkdownFile(filePath, true)
					if err != nil {
						log.Printf("Ошибка при обработке файла %s: %v", fileName, err)
						SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2(fmt.Sprintf("Ошибка при обработке файла %s: %v", fileName, err)))
						continue
					}

					SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2(fmt.Sprintf("📁 Файл: %s", fileName)))
					time.Sleep(500 * time.Millisecond)

					err = SendMarkdownMessage(config.BotToken, chatID, outputMarkdown)
					if err != nil {
						log.Printf("Ошибка при отправке файла %s: %v", fileName, err)
						SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2(fmt.Sprintf("Не удалось отправить содержимое файла %s: %v", fileName, err)))
					} else {
						log.Printf("Файл %s успешно отправлен в чат %d", fileName, chatID)
					}
				} else if messageText == "/help" {
					// Отправляем справку
					helpText := `
Доступные команды:
/start или /files - отправить все тестовые файлы
/file имя_файла - отправить конкретный файл
/help - показать эту справку
					`
					SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2(helpText))
				} else {
					// Для любого другого сообщения отправляем обрабатываем его как Markdown
					// и отправляем обратно в формате Markdown V2
					p := parser.NewWithExtensions(parser.CommonExtensions)
					doc := p.Parse([]byte(messageText))
					renderer := md2.NewRenderer()
					output := Render(doc, renderer)

					// Отправляем обработанное сообщение
					err := SendMarkdownMessage(config.BotToken, chatID, string(output))
					if err != nil {
						log.Printf("Ошибка при отправке обработанного сообщения: %v", err)
						SendMarkdownMessage(config.BotToken, chatID, escapeMarkdownV2(fmt.Sprintf("Ошибка: %v", err)))
					}
				}
			}
		}

		// Небольшая пауза между запросами (чтобы не перегружать API)
		time.Sleep(1 * time.Second)
	}
}
