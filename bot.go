package main

import (
	"cw-arena-watcher/lib"
	"cw-arena-watcher/messages"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"strconv"
	"time"
)

func InitBot(telegramToken string, logger *logrus.Logger, consumer *kafka.Consumer) error {
	logger.Debug("Initializing Telegram bot")
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token:  telegramToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
	if err != nil {
		return err
	}

	bot.Handle("/auth", func(message *telebot.Message) {
		fmt.Println(message.Chat.ID)
	})

	consumer.SubscribeTopics([]string{"cw3-duels"}, nil)

	defer bot.Start()

	chat, err := bot.ChatByID(lib.GetEnv("CWAW_CHANNEL_ID", "-1001451023900"))

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var message messages.DuelMessage
			err = json.Unmarshal([]byte(msg.Value), &message)
			if err != nil {
				sentry.CaptureException(err)
				logger.Error(fmt.Sprintf("Decoder error: %v (%v)\n", err, msg))
			}

			if message.Winner.Tag != "" {
				message.Winner.Tag = "[" + message.Winner.Tag + "]"
			}
			if message.Loser.Tag != "" {
				message.Loser.Tag = "[" + message.Loser.Tag + "]"
			}

			msgString := "Победитель: " +
				message.Winner.Castle +
				message.Winner.Tag +
				message.Winner.Name +
				" 🏅" + strconv.Itoa(message.Winner.Level) +
				" ❤" + strconv.Itoa(message.Winner.Health) + " \n" +
				"Проигравший: " +
				message.Loser.Castle +
				message.Loser.Tag +
				message.Loser.Name +
				" 🏅" + strconv.Itoa(message.Loser.Level) +
				" ❤" + strconv.Itoa(message.Loser.Health)
			
			if message.IsChallenge {
				msgString += "\n" + "<b>Дружеская дуэль</b>"
			}
			
			if message.IsGuildDuel {
				msgString += "\n" + "<b>Гильдейская дуэль</b>"
			}

			_, err = bot.Send(chat, msgString, telebot.ParseMode(telebot.ModeHTML))
			if err != nil {
				sentry.CaptureException(err)
				logger.Error(err)
			}
			logger.Trace(fmt.Sprintf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value)))
		} else {
			sentry.CaptureException(err)
			logger.Error(fmt.Sprintf("Consumer error: %v (%v)\n", err, msg))
		}
	}
}
