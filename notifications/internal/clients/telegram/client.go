package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Client struct {
	bot *tgbotapi.BotAPI
	ChatID int64
}

// Create a client for Telegram
func New(apiKey string, chatId int64) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, err
	}


	return &Client{
		bot: bot,
		ChatID: chatId,
	}, nil
}

// Send message to chat
func (c *Client) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(c.ChatID, message)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)
	return err
}