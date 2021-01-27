package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"lyag/core"
	"lyag/data"
	"lyag/storage"
)

type LyagBot interface {
	ProcessIncomingMessages()
	Run(checkSent bool)
}

func NewLyagBot(settings *core.Settings, yClients core.YClients, storage storage.Checker) (LyagBot, error) {
	telegramAPI, err := tgbotapi.NewBotAPI(settings.TelegramToken)
	if err != nil {
		return nil, err
	}

	telegramAPI.Debug = false

	bot := &lyagBot{
		telegramAPI: telegramAPI,
		yClients:    yClients,
		storage:     storage,
		settings:    settings,
	}

	return bot, nil
}

type lyagBot struct {
	telegramAPI *tgbotapi.BotAPI
	yClients    core.YClients
	storage     storage.Checker
	settings    *core.Settings
}

func needNotify(s storage.Checker, seance *core.Seance) bool {
	wasSent, err := s.WasSent(seance)

	if err != nil {
		log.Printf("Error while checking sent for %s; %s", seance, err)
	}

	if wasSent {
		log.Printf("Already sent for %s, skip", seance)
		return false
	}

	return true
}

func (b *lyagBot) Run(checkSent bool) {

	seances := b.yClients.AvailabilityFor(b.settings.StaffIds, b.settings.DaysDepth)

	for _, staffId := range b.settings.StaffIds {
		staffSeances := make(map[string][]core.Seance)

		for _, seance := range seances {
			if seance.StaffId != staffId {
				continue
			}

			if checkSent && needNotify(b.storage, &seance) == false {
				// prolong status if seance known, set TTL 1+ operation
				ttl := b.settings.Duration + (b.settings.Duration / 3)
				err := b.storage.SetStatus(&seance, true, ttl)
				if err != nil {
					log.Printf("Error while set status for %s; %s", seance, err)
				}
			}

			date := seance.Time.Format("02.01.2006")
			staffSeances[date] = append(staffSeances[date], seance)
		}

		if len(staffSeances) == 0 {
			log.Printf("No seances for %s", staffId)
			continue
		}

		err := b.sendTimes(b.settings.ChatID, staffId, &staffSeances)

		if err != nil {
			log.Printf("Error while sending notification: %s", err)
		}
	}
}

func (b *lyagBot) sendTimes(chatId int64, staffId string, seances *map[string][]core.Seance) error {
	staffName := data.StaffMap[staffId]

	for d, ss := range *seances {
		message := fmt.Sprintf("*%s*: *%s*", staffName, d)

		timeButtons := make([]tgbotapi.InlineKeyboardButton, len(ss))

		for i, s := range ss {
			bData := s.GetKey()
			timeButtons[i] = tgbotapi.NewInlineKeyboardButtonData(s.Time.Format("15:04"), bData)
		}

		msg := tgbotapi.NewMessage(chatId, message)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(timeButtons...))
		msg.ParseMode = "markdown"

		_, err := b.telegramAPI.Send(msg)

		if err != nil {
			return err
		}
	}
	return nil
}

func (b *lyagBot) ProcessIncomingMessages() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := b.telegramAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {

			chatId := update.CallbackQuery.Message.Chat.ID

			if chatId != b.settings.ChatID {
				log.Printf("Unknown chat ID for CallbackQuery %v", chatId)
				continue
			}

			key := update.CallbackQuery.Data
			staffId, date, err := core.ParseKey(key)

			if err != nil {
				log.Printf("Error while parsing key %s: %s", key, err)
				continue
			}

			err = b.yClients.Book(staffId, date)

			if err == nil {
				txt := fmt.Sprintf("Booking created on %v", date)
				log.Printf(txt)
				msg := tgbotapi.NewMessage(chatId, txt)
				_, err = b.telegramAPI.Send(msg)
				if err != nil {
					log.Printf("Error while send confirm message! %s", err)
				}
			} else {
				txt := fmt.Sprintf("Error while booking %s, %s", key, err)
				log.Printf(txt)
				msg := tgbotapi.NewMessage(chatId, txt)
				_, err = b.telegramAPI.Send(msg)
				if err != nil {
					log.Printf("Error while send booking error message! %s", err)
				}
			}
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			switch cmd := update.Message.Command(); cmd {
			case "all":
				b.Run(false)
			}
		}
	}
}
