package main

import (
	"andia/config"
	"andia/download"
	"andia/engine"
	"andia/watermark"
	"fmt"
	"os"
	"time"

	"github.com/tucnak/telebot"
)

func main() {
	config := config.Read()
	botEngine := engine.Engine{}
	bot, err := telebot.NewBot(config.Token)
	if err != nil {
		panic(err)
	}

	messages := make(chan telebot.Message)
	bot.Messages = messages

	go func() {
		for message := range messages {
			// this robot works only for one person
			if message.Sender.ID != config.AdminUser {
				// not admin
				continue
			}

			if message.Text == "START" && botEngine.State == engine.STATE_READY {
				botEngine.State = engine.STATE_WAITING
				bot.SendMessage(message.Sender, "لطفا کپشن مورد نظر خود را بنویسید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						CustomKeyboard: [][]string{
							[]string{
								"STOP",
							},
						},
						ResizeKeyboard: true,
					},
				})
				continue
			}

			if message.Text == "STOP" {
				botEngine.State = engine.STATE_READY
				bot.SendMessage(message.Sender, "غیر فعال شد\n\n"+"برای شروع روی گزینه START کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						CustomKeyboard: [][]string{
							[]string{
								"START",
							},
						},
						ResizeKeyboard: true,
					},
				})
				continue
			}

			if botEngine.State == engine.STATE_WAITING {
				botEngine.Caption = message.Text
				botEngine.State = engine.STATE_LISTENING
				bot.SendMessage(message.Sender, "لطفا تصاویر خود را برای ربات بفرستید تا در کانال منتشر شود ، در صورت اتمام روی گزینه STOP کلیک کنید", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						CustomKeyboard: [][]string{
							[]string{
								"STOP",
							},
						},
						ResizeKeyboard: true,
					},
				})
				continue
			}

			if botEngine.State == engine.STATE_LISTENING && len(message.Photo) > 0 {
				photo := message.Photo[len(message.Photo)-1]
				url, err := bot.GetFileDirectURL(photo.FileID)
				if err != nil {
					continue
				}

				err = download.Download(url, "images/temp.jpg")
				if err != nil {
					continue
				}
				err = watermark.Watermark("images/temp.jpg", "temp.jpg")
				if err != nil {
					continue
				}

				file, err := telebot.NewFile("temp.jpg")
				if err != nil {
					continue
				}

				caption := fmt.Sprintf("%s%s", botEngine.Caption, config.FixedFooter)
				err = bot.SendPhoto(telebot.Chat{Type: "channel", Username: config.Channel}, &telebot.Photo{Caption: caption, File: file}, nil)
				if err != nil {
					continue
				}

				os.Remove("temp.jpg")

				continue
			}

			bot.SendMessage(message.Sender, "دستور شما یافت نشد \n\n"+"برای شروع روی گزینه START کلیک کنید", &telebot.SendOptions{
				ReplyMarkup: telebot.ReplyMarkup{
					CustomKeyboard: [][]string{
						[]string{
							"START",
						},
					},
					ResizeKeyboard: true,
				},
			})
		}
	}()

	bot.Start(time.Second)
}
